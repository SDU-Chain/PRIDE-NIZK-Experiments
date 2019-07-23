import xlrd
import xlsxwriter
import numpy as np
import matplotlib.pyplot as plt
import matplotlib as mpl
import pandas as pd


def status(list):
    x = pd.Series((v[0] for v in list))
    return pd.Series([x.count(), x.min(), x.idxmin(), x.quantile(.25), x.median(),
                      x.quantile(.75), x.mean(), x.max(), x.idxmax(), x.mad(), x.var(),
                      x.std(), x.skew(), x.kurt()], index=['总数', '最小值', '最小值位置', '25%分位数',
                                                           '中位数', '75%分位数', '均值', '最大值', '最大值位数', '平均绝对偏差', '方差', '标准差',
                                                           '偏度', '峰度'])


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
    sheet = workbook.add_worksheet('Worksheet')

    row_id = 0
    sheet.write(row_id, 0, 'filename')
    sheet.write(row_id, 1, '25%')
    sheet.write(row_id, 2, '50%')
    sheet.write(row_id, 3, '75%')
    sheet.write(row_id, 4, '95%')
    sheet.write(row_id, 5, '99%')

    for k in sorted(group_verproof_execution_times.keys()):
        df = pd.DataFrame(group_verproof_execution_times[k])
        row_id += 1
        sheet.write(row_id, 0, k)
        sheet.write(row_id, 1, df.quantile(.25)[0])
        sheet.write(row_id, 2, df.quantile(.50)[0])
        sheet.write(row_id, 3, df.quantile(.75)[0])
        sheet.write(row_id, 4, df.quantile(.95)[0])
        sheet.write(row_id, 5, df.quantile(.99)[0])

    workbook.close()
