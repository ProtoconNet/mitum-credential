import os
import csv
import codecs
import re
from operator import itemgetter

def should_exclude_file(file_path, exclude_files, exclude_folders):
    """파일명이 제외 리스트에 있는지, 또는 파일이 제외할 폴더 내에 있는지 확인합니다."""
    for exclude_file in exclude_files:
        if file_path.endswith(exclude_file):
            return True
    for exclude_folder in exclude_folders:
        if exclude_folder in file_path:
            return True
    return False

def is_process_function(line):
    """'Process' 함수의 시작인지 확인합니다."""
    return line.strip().startswith('func') and ' Process' in line

def is_new_process_function(line):
    return line.strip().startswith('func') and 'types.GetNewProcessor' in line

def find_function_end(lines, start_index):
    """'Process' 함수의 끝을 찾습니다."""
    for i, line in enumerate(lines[start_index:], start_index):
        # 줄의 내용이 오직 '}'만 포함하고 있는지 직접 확인합니다.
        stripped_line = line.rstrip()
        # 이것은 함수의 마지막 블록으로 간주할 수 있습니다.
        if stripped_line == '}':
            return i
    return None  # 함수의 끝을 찾지 못한 경우

def process_line(line):
    """줄의 내용을 콜론으로 나누어 분류문구1, 분류문구2, 그리고 세부 내용을 반환합니다."""
    parts = line.split(':', 2)  # 최대 3개의 부분으로 나눕니다.
    if len(parts) == 3:
        return parts[0].strip(), parts[1].strip(), parts[2].strip()  # 분류문구1, 분류문구2, 세부 내용
    elif len(parts) == 2:
        return parts[0].strip(), '', parts[1].strip()  # 분류문구1, 분류문구2 없음, 세부 내용
    else:
        return '', '', line.strip()  # 분류문구 없음, 세부 내용

def write_csv(filename, rows, header):
    """주어진 헤더와 데이터로 CSV 파일을 작성합니다."""
    with open(filename, 'w', newline='', encoding='utf-8') as csvfile:
        csvwriter = csv.writer(csvfile)
        csvwriter.writerow(header)
        for row in rows:
            csvwriter.writerow(row)

def find_and_process_lines(folder_path, search_strings, exclude_files, exclude_folders, exclude_path_pattern):
    search_pattern = re.compile('|'.join(map(re.escape, search_strings)))
    exclude_pattern = re.compile(exclude_path_pattern)
    quote_content_pattern = re.compile(r'"([^"]*)"')
    data = []

    for root, dirs, files in os.walk(folder_path):
         for file in files:
              file_path = os.path.join(root, file)
              if file.endswith('.go') and not should_exclude_file(file_path, exclude_files, exclude_folders):
                  display_file_path = re.sub(exclude_pattern, '', file_path)
                  try:
                      with codecs.open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
                          lines = f.readlines()
                          i = 0
                          while i < len(lines):
                              line = lines[i]
                              if is_process_function(line):
                                  end_index = find_function_end(lines, i + 1)
                                  if end_index:
                                      i = end_index  # 'Process' 함수의 끝으로 이동
                                      continue
                              if is_new_process_function(line):
                                  end_index = find_function_end(lines, i + 1)
                                  if end_index:
                                      i = end_index  # 'Process' 함수의 끝으로 이동
                                      continue


                              if search_pattern.search(line):
                                  quoted_contents = ' '.join(quote_content_pattern.findall(line))
                                  if quoted_contents:
                                      category1, category2, detail = process_line(quoted_contents)
                                      data.append([display_file_path, i + 1, category1, category2, detail])
                              i += 1
                  except Exception as e:
                      print(f"Error reading file {file_path}: {e}")

    # 파일명에 따라 정렬하여 첫 번째 CSV 파일을 작성
    sorted_by_filename = sorted(data, key=itemgetter(0))
    write_csv('output_by_filename.csv', sorted_by_filename, ['파일명', '줄번호', '분류문구1', '분류문구2', '세부 내용'])

    # 분류문구1에 따라 정렬하여 두 번째 CSV 파일을 작성
    sorted_by_category = sorted(data, key=lambda x: (x[2] == '', x[2]))
    write_csv('output_by_category.csv', sorted_by_category, ['파일명', '줄번호', '분류문구1', '분류문구2', '세부 내용'])

# Golang 프로젝트 폴더 경로
folder_path = '/Users/soonkukkang/go/src/github.com/ProtoconNet/mitum-credential'
# 탐색하고 싶은 문자열들의 리스트
search_strings = ['errors.Wrap', 'errors.Wrapf', 'Errorf', 'base.NewBaseOperationProcessReasonError']
# 탐색에서 제외할 파일명들의 리스트
exclude_files = ['main.go', 'error.go', 'operation_processor.go']
# CSV에 기록될 때 파일명에서 제외할 경로 부분
exclude_path_pattern = '/Users/soonkukkang/go/src/github.com/ProtoconNet/mitum-credential'
# 탐색에서 제외할 폴더 경로 목록
exclude_folders = ['/Users/soonkukkang/go/src/github.com/ProtoconNet/mitum-credential/cmds', '/Users/soonkukkang/go/src/github.com/ProtoconNet/mitum-credential/utils', '/Users/soonkukkang/go/src/github.com/ProtoconNet/mitum-credential/digest']


find_and_process_lines(folder_path, search_strings, exclude_files, exclude_folders, exclude_path_pattern)
