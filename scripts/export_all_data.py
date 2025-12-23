import subprocess
import os
import sys

SCRIPTS = [
    'export_items.py',
    'export_item_sets.py',
    'export_objects_mysql.py',
    'export_quests.py',
    'export_creatures.py',
    'export_factions.py',
    'export_spells.py',
    'export_loot.py',
    'export_metadata.py',
    'extract_atlasloot.py',
]

def run_scripts():
    cwd = os.path.dirname(os.path.abspath(__file__))
    print(f"Starting full data export from {cwd}...")
    
    for script in SCRIPTS:
        script_path = os.path.join(cwd, script)
        if not os.path.exists(script_path):
             print(f"Warning: Script {script} not found, skipping.")
             continue
             
        print(f"--- Running {script} ---")
        try:
            subprocess.run([sys.executable, script], cwd=cwd, check=True)
            print(f"--- {script} Success ---\n")
        except subprocess.CalledProcessError as e:
            print(f"!!! Error running {script}: {e}\n")
            
    print("All exports completed.")

if __name__ == '__main__':
    run_scripts()
