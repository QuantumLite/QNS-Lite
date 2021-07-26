# -*- coding: utf-8 -*-
"""
Created on Thu Apr 29 18:07:24 2021

@author: Amirali
"""

import matplotlib.pyplot as plt
import numpy as np

filePath = "../Data/experiment3.txt"
f = open(filePath, 'r')
try:
    data = f.readlines()
    splitted = data[1].split("[")
    splitted = splitted[1].split("]")
    splitted = splitted[0]
    splitted = splitted.split(", ")
    intSplitted = np.zeros((len(data)-1, len(splitted)), dtype=int)
    for i in range(1, len(data)):
        splitted = data[i].split("[")
        splitted = splitted[1].split("]")
        splitted = splitted[0]
        splitted = splitted.split(", ")
        cntr = 0
        for num in splitted:
            intSplitted[i-1][cntr] = int(num)
            cntr = cntr + 1
    print(data)
    print(data[1])
    print(splitted[0])
    print(type(splitted[0]))
    print(intSplitted)
finally:
    f.close()

MG_NOPP = intSplitted[0]
MG_OPP = intSplitted[1]
NL_NOPP = intSplitted[2]
NL_OPP = intSplitted[3]
QP_NOPP = intSplitted[4]
QP_OPP = intSplitted[5]
sizes = [5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]
improvement_MG = np.divide(np.subtract(MG_NOPP, MG_OPP), MG_NOPP)
improvement_NL = np.divide(np.subtract(NL_NOPP, NL_OPP), NL_NOPP)
improvement_QP = np.divide(np.subtract(QP_NOPP, QP_OPP), QP_NOPP)

plt.plot(sizes, MG_NOPP, "-*", label='MG_NOPP')
plt.plot(sizes, MG_OPP, "-o", label='MG_OPP')
plt.plot(sizes, NL_NOPP, "-*", label='NL_NOPP')
plt.plot(sizes, NL_OPP, "-o", label='NL_OPP')
plt.plot(sizes, QP_NOPP, "-*", label='QP_NOPP')
plt.plot(sizes, QP_OPP, "-o", label='QP_OPP')
plt.legend()
plt.savefig('N20L30PSW0.8PG0.8-algos.pdf', format='pdf')
plt.show()
#plt.subplot(1,2,2)
plt.plot(sizes, improvement_MG, "-*", label='MG')
plt.plot(sizes, improvement_NL, "-o", label='NL')
plt.plot(sizes, improvement_QP, "-o", label='QP')
plt.legend()
plt.savefig('N20L30PSW0.8PG0.8-impros.pdf', format='pdf')

BIGGER_SIZE = 8
AXES_SIZE = 10
MARKER_SIZE = 4
LINE_WIDTH = 0.9

plt.rcParams['pdf.fonttype'] = 42
plt.rcParams['ps.fonttype'] = 42
plt.rcParams["font.family"] = "Times New Roman"
plt.rc('font', size=BIGGER_SIZE)          # controls default text sizes
plt.rc('axes', titlesize=AXES_SIZE)     # fontsize of the axes title
plt.rc('axes', labelsize=AXES_SIZE)    # fontsize of the x and y labels
plt.rc('xtick', labelsize=BIGGER_SIZE)    # fontsize of the tick labels
plt.rc('ytick', labelsize=BIGGER_SIZE)    # fontsize of the tick labels
plt.rc('legend', fontsize=BIGGER_SIZE)

f, ax = plt.subplots(1,2, figsize=(6,2.5), squeeze=False)

plt.xlabel('Size')
ax[0][0].plot(sizes, MG_NOPP, "->", label='MG_NOPP', markersize=MARKER_SIZE, linewidth=LINE_WIDTH)
ax[0][0].plot(sizes, MG_OPP, "-o", label='MG_OPP', markersize=MARKER_SIZE, linewidth=LINE_WIDTH)
ax[0][0].plot(sizes, NL_NOPP, "->", label='NL_NOPP', markersize=MARKER_SIZE, linewidth=LINE_WIDTH)
ax[0][0].plot(sizes, NL_OPP, "-o", label='NL_OPP', markersize=MARKER_SIZE, linewidth=LINE_WIDTH)
ax[0][0].plot(sizes, QP_NOPP, "-*", label='QP_NOPP', markersize=MARKER_SIZE, linewidth=LINE_WIDTH)
ax[0][0].plot(sizes, QP_OPP, "-o", label='QP_OPP', markersize=MARKER_SIZE, linewidth=LINE_WIDTH)
ax[0][0].set_ylabel('Average total waiting time (slots)')
ax[0][0].set_xlabel('Size')
ax[0][0].legend(loc='best', fancybox=True,frameon=False,framealpha=0.8)
plt.ylabel('Improvement ratio')
plt.xlabel('Size')
ax[0][1].plot(sizes, improvement_MG, "->", label='MG', markersize=MARKER_SIZE, linewidth=LINE_WIDTH)
ax[0][1].plot(sizes, improvement_NL, "-o", label='NL', markersize=MARKER_SIZE, linewidth=LINE_WIDTH)
ax[0][1].plot(sizes, improvement_QP, "-o", label='QP', markersize=MARKER_SIZE, linewidth=LINE_WIDTH)
ax[0][1].legend(loc='best', fancybox=True,frameon=False,framealpha=0.8)
plt.subplots_adjust(wspace=0.35)
plt.subplots_adjust(hspace=0.1)
plt.tight_layout()
plt.savefig('N20L30PSW0.8PG0.8-sideBySide.pdf', format='pdf')
plt.show()