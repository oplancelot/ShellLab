#!/usr/bin/env python3
"""
Export GameObject data from MySQL to JSON files
Required for initializing the SQLite database without MySQL connection
"""
import mysql.connector
import json
import os
import db_config


def export_to_json():
    print("=" * 80)
    print("Exporting GameObject data from MySQL to JSON")
    print("=" * 80)
    
    # Ensure data directory exists
    DATA_DIR = os.path.join(os.path.dirname(__file__), '..', 'data')
    os.makedirs(DATA_DIR, exist_ok=True)
    
    
    # 1. Export gameobject_template (tw_world)
    print("\n1. Exporting gameobject_template...")
    try:
        conn = mysql.connector.connect(**db_config.get_mysql_config(db_config.get_tw_world_db()))
        cursor = conn.cursor(dictionary=True)
        
        cursor.execute("""
            SELECT entry, name, type, displayId, size,
                   data0, data1, data2, data3, data4, data5, data6, data7
            FROM gameobject_template
            WHERE name != '' AND name IS NOT NULL
        """)
        
        objects = cursor.fetchall()
        print(f"   Found {len(objects)} objects")
        
        # Mapping to Go struct field names for easier import
        mapped_objects = []
        for obj in objects:
            mapped_objects.append({
                "entry": obj['entry'],
                "name": obj['name'],
                "type": obj['type'],
                "displayId": obj['displayId'],
                "size": float(obj['size']) if obj['size'] is not None else 0.0,
                "data": [
                    obj['data0'] or 0,
                    obj['data1'] or 0,
                    obj['data2'] or 0,
                    obj['data3'] or 0,
                    obj['data4'] or 0,
                    obj['data5'] or 0,
                    obj['data6'] or 0,
                    obj['data7'] or 0
                ]
            })
            
        output_path = os.path.join(DATA_DIR, 'objects.json')
        with open(output_path, 'w', encoding='utf-8') as f:
            json.dump(mapped_objects, f, ensure_ascii=False) # No indent for smaller size
            
        print(f"   ✓ Saved to {output_path}")
        
        conn.close()
    except Exception as e:
        print(f"   ❌ Error: {e}")
        return

    # 2. Export aowow_lock (aowow)
    print("\n2. Exporting aowow_lock...")
    try:
        conn = mysql.connector.connect(**db_config.get_mysql_config(db_config.get_aowow_db()))
        cursor = conn.cursor(dictionary=True)
        
        cursor.execute("SELECT * FROM aowow_lock")
        locks = cursor.fetchall()
        print(f"   Found {len(locks)} lock entries")
        
        output_path = os.path.join(DATA_DIR, 'locks.json')
        with open(output_path, 'w', encoding='utf-8') as f:
            json.dump(locks, f, ensure_ascii=False)
            
        print(f"   ✓ Saved to {output_path}")
        
        conn.close()
    except Exception as e:
        print(f"   ❌ Error: {e}")
        return

    print("\n✓ Export completed successfully!")

if __name__ == "__main__":
    export_to_json()
