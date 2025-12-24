import os
import json
import mysql.connector
import db_config

DATA_DIR = os.path.join(os.path.dirname(os.path.dirname(__file__)), 'data')

def export_icons():
    print("Connecting to MySQL...")
    conn = mysql.connector.connect(**db_config.get_mysql_config())
    cursor = conn.cursor()
    
    print("Fetching item_display_info...")
    cursor.execute("SELECT ID, icon FROM item_display_info")
    
    icons = {}
    for (id, icon) in cursor.fetchall():
        icons[id] = icon
        
    cursor.close()
    conn.close()
    
    out_path = os.path.join(DATA_DIR, 'item_icons.json')
    print(f"Exporting {len(icons)} icons to {out_path}...")
    
    with open(out_path, 'w', encoding='utf-8') as f:
        json.dump(icons, f, ensure_ascii=False)
        
    print("Done!")

if __name__ == "__main__":
    export_icons()
