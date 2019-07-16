#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Fri Mar  1 15:22:13 2019

@author: rkd
"""

import xlrd
import numpy as np
import matplotlib.pyplot as plt
import pandas as pd

workbook = xlrd.open_workbook('exp1_1.xlsx')
sheet = workbook.sheet_by_name(workbook.sheet_names()[0])
# print(sheet.name, sheet.nrows, sheet.ncols)
row_number_of_commitments = sheet.col_values(0)
row_data_first_4_bytes = sheet.col_values(2)
row_transaction_execution_times = sheet.col_values(5)
row_transaction_latency_times = sheet.col_values(8)
boxPlotData = []

group_transaction_execution_times = dict()
group_transaction_latency_times = dict()
for i in range(0, len(row_number_of_commitments)):
    v = 0
    try:
        v = int(row_number_of_commitments[i])
    except ValueError:
        continue

    # VerProof
    if row_data_first_4_bytes[i] == "98d69d92":
        if v in group_transaction_execution_times:
            group_transaction_execution_times[v].append(int(row_transaction_execution_times[i]) / 1000000)
        else:
            group_transaction_execution_times[v] = [int(row_transaction_execution_times[i]) / 1000000]
    # Latency
    elif row_data_first_4_bytes[i]=="027a1f7a":
        if v in group_transaction_latency_times:
            group_transaction_latency_times[v].append(int(row_transaction_latency_times[i]) / 1000000)
        else:
            group_transaction_latency_times[v] = [int(row_transaction_latency_times[i]) / 1000000]

# 注：group_transaction_execution_times.keys() 出现了两次，要求这两次的返回列表顺序必须相同。不知道是否是未定义行为
for k in group_transaction_execution_times.keys():
    boxPlotData.append(group_transaction_execution_times[k])

print(boxPlotData)

plt.style.use("ggplot")
fig = plt.figure(figsize=(20, 12))

plt.boxplot(x=boxPlotData,
            patch_artist=True,
            boxprops={'color': 'black', 'facecolor': '#9999ff', 'linewidth': 1.5},
            flierprops={'marker': 'o', 'markerfacecolor': 'red', 'color': 'black', 'linewidth': 1.5},
            medianprops={'linestyle': '-', 'color': 'blue', 'linewidth': 1.5},
            whiskerprops={'linewidth': 1.5},
            capprops={'linewidth': 1.5},
            labels=group_transaction_execution_times.keys())

# plt.ylim(0, 50)
plt.xlabel('Number of commitments', fontsize=32, color='black', labelpad=20)
plt.ylabel('VerProof (ms)', fontsize=32, color='black', labelpad=20)
plt.xticks(fontsize=32, color='black')
plt.yticks(fontsize=32, color='black')
plt.show()
fig.savefig("graph-1.pdf", bbox_inches="tight")


boxPlotData = []
# 注：group_transaction_execution_times.keys() 出现了两次，要求这两次的返回列表顺序必须相同。不知道是否是未定义行为
for k in group_transaction_latency_times.keys():
    boxPlotData.append(group_transaction_latency_times[k])

print(boxPlotData)

plt.style.use("ggplot")
fig = plt.figure(figsize=(20, 12))

plt.boxplot(x=boxPlotData,
            patch_artist=True,
            boxprops={'color': 'black', 'facecolor': '#9999ff', 'linewidth': 1.5},
            flierprops={'marker': 'o', 'markerfacecolor': 'red', 'color': 'black', 'linewidth': 1.5},
            medianprops={'linestyle': '-', 'color': 'blue', 'linewidth': 1.5},
            whiskerprops={'linewidth': 1.5},
            capprops={'linewidth': 1.5},
            labels=group_transaction_latency_times.keys())

# plt.ylim(0, 50)
plt.xlabel('Number of commitments', fontsize=32, color='black', labelpad=20)
plt.ylabel('VerProof (ms)', fontsize=32, color='black', labelpad=20)
plt.xticks(fontsize=32, color='black')
plt.yticks(fontsize=32, color='black')
plt.show()
fig.savefig("graph-2.pdf", bbox_inches="tight")
