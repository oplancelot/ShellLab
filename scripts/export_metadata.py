import db_config
import mysql.connector
import json
import os

mysql_config = db_config.get_mysql_config(database='aowow')
DATA_DIR = os.path.join(os.path.dirname(__file__), '..', 'data')

def export_metadata():
    print("Connecting to MySQL (aowow)...")
    try:
        conn = mysql.connector.connect(**mysql_config)
        cursor = conn.cursor(dictionary=True)
        
        # 1. Export Zones
        print("Exporting aowow_zones...")
        cursor.execute("SELECT * FROM aowow_zones")
        zones = cursor.fetchall()
        with open(os.path.join(DATA_DIR, 'zones.json'), 'w', encoding='utf-8') as f:
            json.dump(zones, f, ensure_ascii=False)

        # 2. Export Skills
        print("Exporting aowow_skill...")
        cursor.execute("SELECT * FROM aowow_skill")
        skills = cursor.fetchall()
        with open(os.path.join(DATA_DIR, 'skills.json'), 'w', encoding='utf-8') as f:
            json.dump(skills, f, ensure_ascii=False)

        # 3. Export Skill Line Ability
        print("Exporting aowow_skill_line_ability...")
        cursor.execute("SELECT * FROM aowow_skill_line_ability")
        sl_abilities = cursor.fetchall()
        with open(os.path.join(DATA_DIR, 'skill_line_abilities.json'), 'w', encoding='utf-8') as f:
            json.dump(sl_abilities, f, ensure_ascii=False)
            
        print("Metadata export done!")
        cursor.close()
        conn.close()

    except Exception as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    export_metadata()
