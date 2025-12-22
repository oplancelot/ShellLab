# AOWOW Database 本地化实施计划

## 目标

将 aowow.sql 中的所有表复刻到 SQLite，实现完整的本地数据库浏览器。

## aowow.sql 表结构

### 核心数据表

| 表名                    | 描述         | 来源      |
| ----------------------- | ------------ | --------- |
| `aowow_icons`           | 图标数据     | aowow.sql |
| `aowow_itemset`         | 套装数据     | aowow.sql |
| `aowow_itemenchantment` | 附魔数据     | aowow.sql |
| `aowow_factions`        | 阵营数据     | aowow.sql |
| `aowow_factiontemplate` | 阵营模板     | aowow.sql |
| `aowow_zones`           | 区域数据     | aowow.sql |
| `aowow_spellicons`      | 法术图标     | aowow.sql |
| `aowow_spellduration`   | 法术持续时间 | aowow.sql |
| `aowow_spellrange`      | 法术范围     | aowow.sql |
| `aowow_spellradius`     | 法术半径     | aowow.sql |
| `aowow_spellcasttimes`  | 施法时间     | aowow.sql |
| `aowow_spelldispeltype` | 驱散类型     | aowow.sql |
| `aowow_spellmechanic`   | 法术机制     | aowow.sql |
| `aowow_lock`            | 锁定数据     | aowow.sql |
| `aowow_char_titles`     | 称号         | aowow.sql |

### tw_world 数据表（需要导入）

| 表名                  | 描述     | 状态      |
| --------------------- | -------- | --------- |
| `item_template`       | 物品数据 | ✅ 已导入 |
| `creature_template`   | NPC 数据 | ❌ 待导入 |
| `gameobject_template` | 游戏对象 | ❌ 待导入 |
| `quest_template`      | 任务数据 | ❌ 待导入 |
| `spell_template`      | 法术数据 | ❌ 待导入 |

## 实施步骤

### 阶段 1：扩展 SQLite 数据库架构

1. 创建 aowow 相关表的 SQLite 定义
2. 修改 db_import 脚本支持导入这些表

### 阶段 2：数据导入

1. 从 MySQL 导出 aowow 表到 JSON 或直接导入 SQLite
2. 从 tw_world 导出 creature_template, gameobject_template, quest_template, spell_template

### 阶段 3：后端 API

1. 添加 NPC 浏览 API
2. 添加 Quest 浏览 API
3. 添加 Object 浏览 API
4. 添加 Spell 浏览 API
5. 添加 Faction 浏览 API
6. 添加 Item Set 浏览 API

### 阶段 4：前端 Database 页面

1. 重构 Database 页面，添加多个 Tab
2. 实现 Items 三级分类（Class > SubClass > InventoryType）
3. 实现 NPCs 二级分类（Type > Family）
4. 实现 Quests 二级分类（Type > Zone）
5. 实现 Item Sets 浏览
6. 实现 Spells 浏览
7. 实现 Objects 浏览

## 优先级建议

**高优先级：**

1. Item Sets - aowow_itemset 表已有，可直接使用
2. Items 三级分类完善

**中优先级：** 3. NPCs - 需要导入 creature_template 4. Quests - 需要导入 quest_template

**低优先级：** 5. Objects - 需要导入 gameobject_template 6. Spells - 需要导入 spell_template 7. Factions

## 下一步行动

请确认：

1. MySQL 连接信息（用于验证表结构）
2. 优先实现哪些功能？
3. 是否需要先完善现有的 Items 分类？
