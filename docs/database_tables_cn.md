# MySQL 数据库表结构文档

## 数据库概览

| 数据库     | 表数量 | 用途                                         |
| ---------- | ------ | -------------------------------------------- |
| `aowow`    | 22     | AOWOW 网站数据（图标、法术、阵营等元数据）   |
| `tw_world` | 180+   | Turtle WoW 游戏世界数据（物品、NPC、任务等） |

---

## aowow 数据库（22 表）

### 角色相关

| 表名                | 用途     |
| ------------------- | -------- |
| `aowow_char_titles` | 角色称号 |

### 阵营相关

| 表名                    | 用途                          |
| ----------------------- | ----------------------------- |
| `aowow_factions`        | 阵营数据                      |
| `aowow_factiontemplate` | 阵营模板（包含友好/敌对关系） |

### 物品相关

| 表名                    | 用途                               |
| ----------------------- | ---------------------------------- |
| `aowow_icons`           | 物品/法术图标名称映射              |
| `aowow_itemenchantment` | 物品附魔效果                       |
| `aowow_itemset`         | 套装数据（包含套装部件和套装奖励） |

### 技能相关

| 表名                       | 用途           |
| -------------------------- | -------------- |
| `aowow_skill`              | 技能数据       |
| `aowow_skill_line_ability` | 技能线能力关联 |

### 法术相关

| 表名                    | 用途                       |
| ----------------------- | -------------------------- |
| `aowow_spell`           | 法术基础数据               |
| `aowow_spellcasttimes`  | 施法时间                   |
| `aowow_spelldispeltype` | 驱散类型                   |
| `aowow_spellduration`   | 法术持续时间               |
| `aowow_spellicons`      | 法术图标                   |
| `aowow_spellmechanic`   | 法术机制（如昏迷、恐惧等） |
| `aowow_spellradius`     | 法术半径                   |
| `aowow_spellrange`      | 法术距离                   |

### 其他

| 表名                   | 用途               |
| ---------------------- | ------------------ |
| `aowow_comments`       | 用户评论           |
| `aowow_comments_rates` | 评论评分           |
| `aowow_lock`           | 锁数据（开锁相关） |
| `aowow_news`           | 新闻               |
| `aowow_resistances`    | 抗性数据           |
| `aowow_zones`          | 区域/地图数据      |

---

## tw_world 数据库（主要表）

### 物品相关

| 表名                         | 用途                           |
| ---------------------------- | ------------------------------ |
| `item_template`              | 物品模板（所有物品的基础数据） |
| `item_display_info`          | 物品显示信息                   |
| `item_enchantment_template`  | 物品附魔模板                   |
| `item_loot_template`         | 物品掉落模板                   |
| `item_required_target`       | 物品使用目标要求               |
| `item_transmogrify_template` | 幻化模板                       |
| `locales_item`               | 物品本地化（多语言）           |

### NPC/生物相关

| 表名                          | 用途                         |
| ----------------------------- | ---------------------------- |
| `creature_template`           | NPC/生物模板                 |
| `creature`                    | 生物实例（世界中的位置）     |
| `creature_addon`              | 生物附加数据（光环、武器等） |
| `creature_ai_events`          | 生物 AI 事件                 |
| `creature_ai_scripts`         | 生物 AI 脚本                 |
| `creature_equip_template`     | 生物装备模板                 |
| `creature_groups`             | 生物组                       |
| `creature_involvedrelation`   | 生物任务关联（任务完成）     |
| `creature_loot_template`      | 生物掉落模板                 |
| `creature_movement`           | 生物移动路径                 |
| `creature_movement_template`  | 生物移动模板                 |
| `creature_onkill_reputation`  | 击杀声望奖励                 |
| `creature_questrelation`      | 生物任务关联（任务接取）     |
| `creature_spells`             | 生物法术                     |
| `creature_display_info_addon` | 生物显示信息附加             |
| `locales_creature`            | 生物本地化                   |

### 任务相关

| 表名                   | 用途         |
| ---------------------- | ------------ |
| `quest_template`       | 任务模板     |
| `quest_cast_objective` | 任务施法目标 |
| `quest_end_scripts`    | 任务结束脚本 |
| `quest_greeting`       | 任务问候语   |
| `quest_start_scripts`  | 任务开始脚本 |
| `locales_quest`        | 任务本地化   |

### 游戏对象相关

| 表名                            | 用途                         |
| ------------------------------- | ---------------------------- |
| `gameobject_template`           | 游戏对象模板                 |
| `gameobject`                    | 游戏对象实例（世界中的位置） |
| `gameobject_involvedrelation`   | 游戏对象任务关联（任务完成） |
| `gameobject_loot_template`      | 游戏对象掉落模板             |
| `gameobject_questrelation`      | 游戏对象任务关联（任务接取） |
| `gameobject_scripts`            | 游戏对象脚本                 |
| `gameobject_display_info_addon` | 游戏对象显示信息             |
| `locales_gameobject`            | 游戏对象本地化               |

### 法术相关

| 表名                    | 用途                   |
| ----------------------- | ---------------------- |
| `spell_template`        | 法术模板               |
| `spell_affect`          | 法术影响               |
| `spell_area`            | 区域法术               |
| `spell_chain`           | 法术链（技能升级关系） |
| `spell_disabled`        | 禁用的法术             |
| `spell_effect_mod`      | 法术效果修正           |
| `spell_elixir`          | 药剂类型               |
| `spell_group`           | 法术分组               |
| `spell_learn_spell`     | 学习法术关联           |
| `spell_mod`             | 法术修正               |
| `spell_proc_event`      | 法术触发事件           |
| `spell_scripts`         | 法术脚本               |
| `spell_target_position` | 法术目标位置           |
| `locales_spell`         | 法术本地化             |

### 掉落模板

| 表名                          | 用途               |
| ----------------------------- | ------------------ |
| `creature_loot_template`      | 生物掉落           |
| `gameobject_loot_template`    | 游戏对象掉落       |
| `item_loot_template`          | 物品掉落（如箱子） |
| `disenchant_loot_template`    | 分解掉落           |
| `fishing_loot_template`       | 钓鱼掉落           |
| `mail_loot_template`          | 邮件掉落           |
| `pickpocketing_loot_template` | 扒窃掉落           |
| `reference_loot_template`     | 引用掉落模板       |
| `skinning_loot_template`      | 剥皮掉落           |

### NPC 交互

| 表名                   | 用途         |
| ---------------------- | ------------ |
| `npc_gossip`           | NPC 对话     |
| `npc_text`             | NPC 文本     |
| `npc_trainer`          | NPC 训练师   |
| `npc_trainer_template` | 训练师模板   |
| `npc_vendor`           | NPC 商人     |
| `npc_vendor_template`  | 商人模板     |
| `gossip_menu`          | 对话菜单     |
| `gossip_menu_option`   | 对话菜单选项 |
| `gossip_scripts`       | 对话脚本     |

### 地图/区域相关

| 表名                     | 用途             |
| ------------------------ | ---------------- |
| `area_template`          | 区域模板         |
| `areatrigger_teleport`   | 区域触发器传送点 |
| `areatrigger_template`   | 区域触发器模板   |
| `areatrigger_tavern`     | 旅馆区域触发器   |
| `map_template`           | 地图模板         |
| `game_graveyard_zone`    | 墓地区域         |
| `world_safe_locs_facing` | 世界安全位置朝向 |
| `locales_area`           | 区域本地化       |

### 战场相关

| 表名                      | 用途           |
| ------------------------- | -------------- |
| `battleground_events`     | 战场事件       |
| `battleground_template`   | 战场模板       |
| `battlemaster_entry`      | 战场管理员入口 |
| `creature_battleground`   | 战场生物       |
| `gameobject_battleground` | 战场游戏对象   |

### 阵营/声望

| 表名                            | 用途         |
| ------------------------------- | ------------ |
| `faction`                       | 阵营         |
| `faction_template`              | 阵营模板     |
| `reputation_reward_rate`        | 声望奖励倍率 |
| `reputation_spillover_template` | 声望溢出模板 |
| `locales_faction`               | 阵营本地化   |

### 玩家相关

| 表名                      | 用途             |
| ------------------------- | ---------------- |
| `player_classlevelstats`  | 玩家职业等级属性 |
| `player_levelstats`       | 玩家等级属性     |
| `player_xp_for_level`     | 升级经验值       |
| `playercreateinfo`        | 玩家创建信息     |
| `playercreateinfo_action` | 玩家创建动作栏   |
| `playercreateinfo_item`   | 玩家创建物品     |
| `playercreateinfo_spell`  | 玩家创建法术     |
| `player_factionchange_*`  | 阵营转换相关     |

### 技能相关

| 表名                       | 用途         |
| -------------------------- | ------------ |
| `skill_fishing_base_level` | 钓鱼基础等级 |
| `skill_line_ability`       | 技能线能力   |

### 宠物相关

| 表名                  | 用途         |
| --------------------- | ------------ |
| `pet_levelstats`      | 宠物等级属性 |
| `pet_name_generation` | 宠物名称生成 |
| `pet_spell_data`      | 宠物法术数据 |
| `petcreateinfo_spell` | 宠物创建法术 |
| `collection_pet`      | 宠物收藏     |
| `collection_mount`    | 坐骑收藏     |

### 交通/传送

| 表名                    | 用途         |
| ----------------------- | ------------ |
| `taxi_nodes`            | 飞行点       |
| `taxi_path_transitions` | 飞行路径过渡 |
| `transports`            | 传送门/船    |
| `game_tele`             | 游戏传送点   |
| `locales_taxi_node`     | 飞行点本地化 |

### 脚本相关

| 表名                   | 用途           |
| ---------------------- | -------------- |
| `event_scripts`        | 事件脚本       |
| `generic_scripts`      | 通用脚本       |
| `script_texts`         | 脚本文本       |
| `script_waypoint`      | 脚本路径点     |
| `scripted_areatrigger` | 脚本区域触发器 |
| `scripted_event_id`    | 脚本事件 ID    |

### 游戏事件

| 表名                       | 用途             |
| -------------------------- | ---------------- |
| `game_event_creature`      | 游戏事件生物     |
| `game_event_creature_data` | 游戏事件生物数据 |
| `game_event_gameobject`    | 游戏事件游戏对象 |
| `game_event_mail`          | 游戏事件邮件     |
| `game_event_quest`         | 游戏事件任务     |

### 刷怪池

| 表名                       | 用途           |
| -------------------------- | -------------- |
| `pool_creature`            | 生物池         |
| `pool_creature_template`   | 生物池模板     |
| `pool_gameobject`          | 游戏对象池     |
| `pool_gameobject_template` | 游戏对象池模板 |
| `pool_pool`                | 池的池         |
| `pool_template`            | 池模板         |

### 商城相关

| 表名              | 用途     |
| ----------------- | -------- |
| `shop_categories` | 商城分类 |
| `shop_items`      | 商城物品 |

### 其他

| 表名                 | 用途          |
| -------------------- | ------------- |
| `autobroadcast`      | 自动广播      |
| `broadcast_text`     | 广播文本      |
| `conditions`         | 条件系统      |
| `exploration_basexp` | 探索基础经验  |
| `game_weather`       | 天气          |
| `mangos_string`      | MaNGOS 字符串 |
| `page_text`          | 书页文本      |
| `points_of_interest` | 兴趣点        |
| `reserved_name`      | 保留名称      |
| `sound_entries`      | 声音入口      |
| `variables`          | 变量          |
| `warden_checks`      | 反作弊检查    |
| `warden_scans`       | 反作弊扫描    |

---

## ShellLab 数据库架构

本项目的数据库采用 ETL (Extract-Transform-Load) 流程，从 source (MySQL) 导出到 JSON，再导入到 local SQLite。

### ETL 流程状态

| 模块          | 源表 (MySQL)               | 中间文件 (JSON)        | 目标表 (SQLite)   | 状态    | 导入脚本                  |
| ------------- | -------------------------- | ---------------------- | ----------------- | ------- | ------------------------- |
| **Items**     | `item_template`            | `item_template.json`   | `items`           | ✅ 完成 | `db_import/main.go`       |
| **Objects**   | `gameobject_template`      | `objects.json`         | `objects`         | ✅ 完成 | `export_objects_mysql.py` |
| **Locks**     | `aowow_lock`               | `locks.json`           | `locks`           | ✅ 完成 | `export_objects_mysql.py` |
| **Quests**    | `quest_template`           | `quests.json`          | `quests`          | ✅ 完成 | `export_quests.py`        |
| **Creatures** | `creature_template`        | `creatures.json`       | `creatures`       | ✅ 完成 | `export_creatures.py`     |
| **Factions**  | `aowow_factions`           | `factions.json`        | `factions`        | ✅ 完成 | `export_factions.py`      |
| **Spells**    | `spell_template`           | `spells.json`          | `spells`          | ✅ 完成 | `export_spells.py`        |
| **Loot**      | `creature_loot_template`   | `creature_loot.json`   | `creature_loot`   | ✅ 完成 | `export_loot.py`          |
|               | `reference_loot_template`  | `reference_loot.json`  | `reference_loot`  | ✅ 完成 | `export_loot.py`          |
|               | `gameobject_loot_template` | `gameobject_loot.json` | `gameobject_loot` | ✅ 完成 | `export_loot.py`          |
|               | `item_loot_template`       | `item_loot.json`       | `item_loot`       | ✅ 完成 | `export_loot.py`          |
|               | `disenchant_loot_template` | `disenchant_loot.json` | `disenchant_loot` | ✅ 完成 | `export_loot.py`          |

### SQLite 表结构定义

以下是 ShellLab 本地 SQLite (`data/shelllab.db`) 的主要表结构。

#### 1. Objects (游戏对象)

```sql
CREATE TABLE objects (
    entry INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    type INTEGER NOT NULL,
    display_id INTEGER,
    size REAL,
    data0-7 INTEGER  -- 包含锁ID等关键数据
);
```

#### 2. Locks (锁数据)

用于推导 Objects 的特殊分类（如草药、矿脉）。

```sql
CREATE TABLE locks (
    id INTEGER PRIMARY KEY,
    type1-5 INTEGER,
    prop1-5 INTEGER, -- 关键属性 (1=开锁, 2=草药, 3=采矿)
    req1-5 INTEGER
);
```

#### 3. Quests (任务)

```sql
CREATE TABLE quests (
    entry INTEGER PRIMARY KEY,
    title TEXT,
    min_level INTEGER,
    quest_level INTEGER,
    ... -- 包含详细的任务文本、奖励、目标等
);
```

#### 4. Creatures (生物)

```sql
CREATE TABLE creatures (
    entry INTEGER PRIMARY KEY,
    name TEXT,
    subname TEXT,
    level_min/max INTEGER,
    health_min/max INTEGER,
    creature_type INTEGER,
    creature_rank INTEGER,
    loot_id INTEGER,
    skin_loot_id INTEGER,
    pickpocket_loot_id INTEGER
    ...
);
```

#### 5. Loot Tables (掉落表)

所有掉落表结构相同：

```sql
CREATE TABLE *_loot (
    entry INTEGER,          -- 关联 ID (如 creature.loot_id)
    item INTEGER,           -- 物品 ID (关联 items.entry)
    chance REAL,            -- 掉率
    groupid INTEGER,
    mincount_or_ref INTEGER, -- 正数=数量, 负数=引用其他 loot 表
    maxcount INTEGER
);
```

#### 6. Spells (法术)

用于显示法术详情和计算套装/物品效果。

```sql
CREATE TABLE spells (
    entry INTEGER PRIMARY KEY,
    name TEXT,
    description TEXT,
    effect_base_points1-3 INTEGER,
    effect_die_sides1-3 INTEGER
);
```

#### 7. Factions (阵营)

```sql
CREATE TABLE factions (
    id INTEGER PRIMARY KEY,
    name TEXT,
    description TEXT,
    side INTEGER,      -- 联盟/部落/中立
    category_id INTEGER
);
```
