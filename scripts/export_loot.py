import db_config
import mysql.connector
import json
import os

mysql_config = db_config.get_mysql_config()
DATA_DIR = os.path.join(os.path.dirname(__file__), '..', 'data')

def export_table_to_json(cursor, table_name, json_filename):
    print(f"Fetching {table_name}...")
    try:
        # Check if table exists first (or we can just try query)
        cursor.execute(f"SELECT entry, item, ChanceOrQuestChance, groupid, mincountOrRef, maxcount FROM {table_name}")
        rows = cursor.fetchall()
        print(f"Found {len(rows)} entries in {table_name}")
        
        filepath = os.path.join(DATA_DIR, json_filename)
        with open(filepath, 'w', encoding='utf-8') as f:
            # Optimize: write line by line or use compact JSON? Compact is better for bulk.
            # But converting row dict keys to simple camelCase or keep as is?
            # Let's keep keys simple for Go parsing: 
            # entry, item, chance, groupId, minCountOrRef, maxCount
            compact_rows = []
            for r in rows:
                compact_rows.append({
                    "entry": r['entry'],
                    "item": r['item'],
                    "chance": r['ChanceOrQuestChance'],
                    "groupId": r['groupid'],
                    "minCountOrRef": r['mincountOrRef'],
                    "maxCount": r['maxcount']
                })
            json.dump(compact_rows, f, ensure_ascii=False)
        print(f"Saved to {filepath}")
    except mysql.connector.Error as err:
        print(f"Warning: Could not export {table_name}: {err}")

def export_loot():
    try:
        os.makedirs(DATA_DIR, exist_ok=True)
        
        mysql_conn = mysql.connector.connect(**mysql_config)
        mysql_cursor = mysql_conn.cursor(dictionary=True)
        
        # 1. creature_loot_template
        export_table_to_json(mysql_cursor, "creature_loot_template", "creature_loot_template.json")
        
        # 2. reference_loot_template
        export_table_to_json(mysql_cursor, "reference_loot_template", "reference_loot_template.json")
        
        # 3. gameobject_loot_template
        export_table_to_json(mysql_cursor, "gameobject_loot_template", "gameobject_loot_template.json")
        
        # 4. item_loot_template (disenchanting/milling/prospecting usually)
        export_table_to_json(mysql_cursor, "item_loot_template", "item_loot_template.json")
        
        # 5. disenchant_loot_template
        export_table_to_json(mysql_cursor, "disenchant_loot_template", "disenchant_loot_template.json")

        mysql_cursor.close()
        mysql_conn.close()
        print("Done!")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    export_loot()
