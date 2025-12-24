import db_config
import mysql.connector
import json
import os
import decimal
from datetime import date, datetime

mysql_config = db_config.get_mysql_config()
DATA_DIR = os.path.join(os.path.dirname(__file__), '..', 'data')

# Helper to handle Decimal and Date serialization
class CustomEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, decimal.Decimal):
            return float(obj)
        if isinstance(obj, (date, datetime)):
            return obj.isoformat()
        return super(CustomEncoder, self).default(obj)

def export_items():
    print("Connecting to MySQL...")
    try:
        mysql_conn = mysql.connector.connect(**mysql_config)
        mysql_cursor = mysql_conn.cursor(dictionary=True)
        
        print("Fetching item_template from MySQL...")
        # Select all columns
        mysql_cursor.execute("SELECT * FROM item_template ORDER BY entry")
        
        # Fetch all results
        # Assuming memory is sufficient (usually ~100MB for vanilla/t-wow databases)
        items = mysql_cursor.fetchall()
        print(f"Found {len(items)} items")
        
        mysql_cursor.close()
        mysql_conn.close()
        
        os.makedirs(DATA_DIR, exist_ok=True)
        json_path = os.path.join(DATA_DIR, 'item_template.json')
        print(f"Exporting to {json_path}...")
        
        # Normalize column names to lowercase keys for consistency if needed, 
        # but MySQL connector dictionary cursor usually returns column names as defined in DB.
        # We will keep them as is.
        
        with open(json_path, 'w', encoding='utf-8') as f:
            json.dump(items, f, cls=CustomEncoder, ensure_ascii=False)
            
        print("Done!")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    export_items()
