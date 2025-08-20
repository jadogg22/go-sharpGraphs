# Eastbound Round Trip Revenue Analysis

This script analyzes round-trip freight data, filtering for trips that start or return to Utah. It identifies round trips, computes revenue per mile, and generates visual reports and CSV summaries.

---

## Features

- Filters for eastbound round trips involving Utah
- Calculates outbound/inbound revenue per mile
- Groups by destination state and outputs a scatter plot
- Saves CSV summaries and PNG visualizations in a `reports/` folder

---

## Setup Instructions

### 1. Clone or Download

Clone the repo, or download and unzip it.

---

## 2. Create and Activate a Virtual Environment

### 🐧 macOS/Linux

```bash
python3 -m venv venv
source venv/bin/activate
```

### Windows

```bash
python -m venv venv
venv\Scripts\activate.bat
```

### Windows PowerShell

```bash
python -m venv venv
venv\Scripts\Activate.ps1
```

### installing dependencies

```bash
pip install -r requirements.txt
```

### running the script

```bash
python main.py
```

## final notes

make sure to change the filename of the CSV file in the script to match your data file. its currently "june_2025.csv"
