import xlrd
import xlsxwriter
import numpy as np
import matplotlib.pyplot as plt
import matplotlib as mpl
import pandas as pd

if __name__ == "__main__":

    # Read Excel
    workbook = xlrd.open_workbook('input.xlsx')
    sheet = workbook.sheet_by_name(workbook.sheet_names()[0])

    # print(sheet.name, sheet.nrows, sheet.ncols)
    row_number_of_commitments = sheet.col_values(0)
    row_data_first_4_bytes = sheet.col_values(2)
    row_transaction_execution_times = sheet.col_values(5)
    row_transaction_latency_times = sheet.col_values(8)

    group_verproof_execution_times = dict()
    group_verproof_latency_times = dict()
    group_newsession_execution_times = dict()
    group_newsession_latency_times = dict()

    for i in range(0, len(row_number_of_commitments)):
        v = 0
        try:
            v = int(row_number_of_commitments[i])
        except ValueError:
            continue

        # VerProof
        if row_data_first_4_bytes[i] == "98d69d92":
            if v in group_verproof_execution_times:
                group_verproof_execution_times[v].append(int(row_transaction_execution_times[i]) / 1000000)
                group_verproof_latency_times[v].append(int(row_transaction_latency_times[i]) / 1000000)
            else:
                group_verproof_execution_times[v] = [int(row_transaction_execution_times[i]) / 1000000]
                group_verproof_latency_times[v] = [(int(row_transaction_latency_times[i]) / 1000000)]

        # new session
        elif row_data_first_4_bytes[i] == "027a1f7a":
            if v in group_newsession_execution_times:
                group_newsession_execution_times[v].append(int(row_transaction_execution_times[i]) / 1000000)
                group_newsession_latency_times[v].append(int(row_transaction_latency_times[i]) / 1000000)
            else:
                group_newsession_execution_times[v] = [int(row_transaction_execution_times[i]) / 1000000]
                group_newsession_latency_times[v] = [(int(row_transaction_latency_times[i]) / 1000000)]

    # write Excel

    workbook = xlsxwriter.Workbook('output.xlsx')


    def analyze(dict_list, title):
        sheet = workbook.add_worksheet(title)

        row_id = 0
        col_id = 0

        sheet.write(row_id, col_id, 'filename')
        col_id += 1
        sheet.write(row_id, col_id, '25%')
        col_id += 1
        sheet.write(row_id, col_id, '50%')
        col_id += 1
        sheet.write(row_id, col_id, '85%')
        col_id += 1
        sheet.write(row_id, col_id, '86%')
        col_id += 1
        sheet.write(row_id, col_id, '87%')
        col_id += 1
        sheet.write(row_id, col_id, '88%')
        col_id += 1
        sheet.write(row_id, col_id, '89%')
        col_id += 1
        sheet.write(row_id, col_id, '90%')
        col_id += 1
        sheet.write(row_id, col_id, 'average')
        col_id += 1

        for k in sorted(dict_list.keys()):
            df = pd.DataFrame(dict_list[k])
            row_id += 1
            col_id = 0
            sheet.write(row_id, col_id, k)
            col_id += 1
            sheet.write(row_id, col_id, df.quantile(.25)[0])
            col_id += 1
            sheet.write(row_id, col_id, df.quantile(.50)[0])
            col_id += 1
            sheet.write(row_id, col_id, df.quantile(.85)[0])
            col_id += 1
            sheet.write(row_id, col_id, df.quantile(.86)[0])
            col_id += 1
            sheet.write(row_id, col_id, df.quantile(.87)[0])
            col_id += 1
            sheet.write(row_id, col_id, df.quantile(.88)[0])
            col_id += 1
            sheet.write(row_id, col_id, df.quantile(.89)[0])
            col_id += 1
            sheet.write(row_id, col_id, df.quantile(.90)[0])
            col_id += 1
            sheet.write(row_id, col_id, df.mean()[0])
            col_id += 1


    analyze(group_verproof_execution_times, 'verproof')
    analyze(group_verproof_latency_times, 'prooflatency')

    workbook.close()
