import db_config
import mysql.connector
import json
import os

mysql_config = db_config.get_mysql_config()
DATA_DIR = os.path.join(os.path.dirname(__file__), '..', 'data')

def export_spells():
    print("Connecting to MySQL...")
    try:
        mysql_conn = mysql.connector.connect(**mysql_config)
        mysql_cursor = mysql_conn.cursor(dictionary=True)
        
        print("Fetching spells from MySQL...")
        # Note: spell_template usually in tw_world (or world db)
        mysql_cursor.execute("""
            SELECT 
                entry, name, description,
                effectBasePoints1, effectBasePoints2, effectBasePoints3,
                effectDieSides1, effectDieSides2, effectDieSides3
            FROM spell_template
            ORDER BY entry
        """)
        spells = mysql_cursor.fetchall()
        print(f"Found {len(spells)} spells")
        mysql_cursor.close()
        mysql_conn.close()
        
        os.makedirs(DATA_DIR, exist_ok=True)
        json_path = os.path.join(DATA_DIR, 'spells.json')
        print(f"Exporting to {json_path}...")
        
        with open(json_path, 'w', encoding='utf-8') as f:
            json.dump(spells, f, ensure_ascii=False)
            
        print("Done!")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    export_spells()
