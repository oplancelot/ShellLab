import db_config
import mysql.connector
import json
import os

# aowow_itemset is in 'aowow' database
mysql_config = db_config.get_mysql_config(database='aowow')
DATA_DIR = os.path.join(os.path.dirname(__file__), '..', 'data')

def export_item_sets():
    print("Connecting to MySQL (aowow)...")
    try:
        mysql_conn = mysql.connector.connect(**mysql_config)
        mysql_cursor = mysql_conn.cursor(dictionary=True)
        
        print("Fetching aowow_itemset...")
        # Select all columns
        mysql_cursor.execute("SELECT * FROM aowow_itemset")
        
        sets = mysql_cursor.fetchall()
        print(f"Found {len(sets)} item sets")
        
        mysql_cursor.close()
        mysql_conn.close()
        
        os.makedirs(DATA_DIR, exist_ok=True)
        json_path = os.path.join(DATA_DIR, 'item_sets.json')
        print(f"Exporting to {json_path}...")
        
        with open(json_path, 'w', encoding='utf-8') as f:
            json.dump(sets, f, ensure_ascii=False)
            
        print("Done!")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    export_item_sets()
