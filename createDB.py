import sqlite3
import requests

url = 'http://192.168.0.62:5000/Logistics/get_yearly_revenue'
url2 = 'http://192.168.0.62:5000/Transportation/get_yearly_revenue'

dbdata = {}

def main():
    conn = sqlite3.connect("Production.db")
    cursor = conn.cursor()

 # Define the SQL statement to create the table
    create_table_sql = '''
        CREATE TABLE IF NOT EXISTS trans_year_rev (
            ID INTEGER PRIMARY KEY AUTOINCREMENT,
            Year INTEGER NOT NULL,
            Week INTEGER NOT NULL,
            TotalRevenue REAL NOT NULL
        );
    '''
    response = requests.get(url2)
    if response.status_code == 200:
        dbdata = response.json()
    print(dbdata)
    
    # Execute the SQL statement to create the table
    cursor.execute(create_table_sql)
    
    # Commit the changes and close the connection
    conn.commit()

    for arr in dbdata["Data"]:
        week = 0
        rev_2021 = 0.0
        rev_2022 = 0.0
        rev_2023 = 0.0


        for key in arr.keys():

            if key == "Name":
                name = arr["Name"]
                week = int(name[1:])
            if key[-3:] != "Act":
                #skip this section
                continue
            if key[:4] == "2021":
                rev_2021 = float(arr[key])
            if key[:4] == "2022":
                rev_2022 = float(arr[key])
            if key[:4] == "2023":
                rev_2023 = float(arr[key])

        # now that we have finished the array lets add the values to the db 
        print(week)
        print(rev_2021)
        print(rev_2022)
        print(rev_2023)
        insert_data(cursor, week, 2021, rev_2021)
        insert_data(cursor, week, 2022, rev_2022)
        insert_data(cursor, week, 2023, rev_2023)
    conn.commit()
    conn.close()

def insert_data(cursor, week, year, revenue):
    sql = """
    INSERT INTO trans_year_rev (Week, Year, TotalRevenue)
    VALUES (?, ?, ?)
    """

    cursor.execute(sql, (week, year, revenue))


main()



