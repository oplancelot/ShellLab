import db_config
import mysql.connector
import json
import os

mysql_config = db_config.get_mysql_config()
DATA_DIR = os.path.join(os.path.dirname(__file__), '..', 'data')

def export_creatures():
    print("Connecting to MySQL...")
    try:
        mysql_conn = mysql.connector.connect(**mysql_config)
        mysql_cursor = mysql_conn.cursor(dictionary=True)
        
        print("Fetching creatures...")
        mysql_cursor.execute("SELECT * FROM creature_template")
        creatures = mysql_cursor.fetchall()
        print(f"Found {len(creatures)} creatures")
        mysql_cursor.close()
        mysql_conn.close()
        
        # Ensure data directory exists
        os.makedirs(DATA_DIR, exist_ok=True)
        
        json_path = os.path.join(DATA_DIR, 'creature_template.json')
        print(f"Exporting to {json_path}...")
        
        with open(json_path, 'w', encoding='utf-8') as f:
            json.dump(creatures, f, ensure_ascii=False)
            
        print("Done!")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    export_creatures()
