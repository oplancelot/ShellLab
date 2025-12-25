import db_config
import mysql.connector
import json
import os

mysql_config = db_config.get_mysql_config()
DATA_DIR = os.path.join(os.path.dirname(__file__), '..', 'data')

def export_quests():
    print("Connecting to MySQL...")
    try:
        mysql_conn = mysql.connector.connect(**mysql_config)
        mysql_cursor = mysql_conn.cursor(dictionary=True)
        
        print("Fetching quests...")
        # Fetching all details needed for the detail page
        query = "SELECT * FROM quest_template"
        mysql_cursor.execute(query)
        quests = mysql_cursor.fetchall()
        print(f"Found {len(quests)} quests")
        mysql_cursor.close()
        mysql_conn.close()
        
        # Ensure data directory exists
        os.makedirs(DATA_DIR, exist_ok=True)
        
        # Export to JSON
        json_path = os.path.join(DATA_DIR, 'quest_template.json')
        print(f"Exporting to {json_path}...")
        
        with open(json_path, 'w', encoding='utf-8') as f:
            # We dump the dictionary list directly. 
            # Go struct tags will need to match these keys (e.g. Title, MinLevel...) 
            # or we map them here. For simplicity, let's keep MySQL column names which are mostly PascalCase.
            json.dump(quests, f, ensure_ascii=False)
            
        print("Done!")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    export_quests()
