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

workbook = xlrd.open_workbook(r'/home/rkd/Desktop/proof-time-ms.xlsx')
sheet = workbook.sheet_by_name(workbook.sheet_names()[0])
# print(sheet.name, sheet.nrows, sheet.ncols)
# name = sheet.col_values(0)
d = sheet.col_values(2)
print(len(d))
boxPlotData = []
for i in range(7):
    tmp = d[20 * i: 20 * (i + 1)]
    boxPlotData.append(tmp)
# print(boxPlotData)
plt.style.use("ggplot")
fig = plt.figure(figsize=(20, 12))

plt.boxplot(x=boxPlotData,
            patch_artist=True,
            boxprops={'color': 'black', 'facecolor': '#9999ff', 'linewidth': 1.5},
            flierprops={'marker': 'o', 'markerfacecolor': 'red', 'color': 'black', 'linewidth': 1.5},
            medianprops={'linestyle': '-', 'color': 'blue', 'linewidth': 1.5},
            whiskerprops={'linewidth': 1.5},
            capprops={'linewidth': 1.5},
            labels=['1', '100', '500', '1000', '5000', '10000', '50000'])
plt.ylim(0, 50)
plt.xlabel('Number of commitments', fontsize=25, color='black', labelpad=20)
plt.ylabel('Time cost (ms)', fontsize=25, color='black', labelpad=20)
plt.xticks(fontsize=25, color='black')
plt.yticks(fontsize=25, color='black')
plt.show()
fig.savefig("/home/rkd/Desktop/proof-time-ms.pdf", bbox_inches="tight")
