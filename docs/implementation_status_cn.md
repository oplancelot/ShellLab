# 本地数据库与功能实现状态

本文档记录 ShellLab 本地数据库（ETL 流程）及前端功能的实现进度。

## 数据库表导入状态 (ETL)

### 核心数据表 (MySQL -> JSON -> SQLite)

| 表名          | 描述     | 状态      | 对应脚本                  |
| ------------- | -------- | --------- | ------------------------- |
| **Items**     | 物品     | ✅ 已完成 | `db_import/main.go`       |
| **Icons**     | 物品图标 | ✅ 已完成 | `import_icons/main.go`    |
| **Objects**   | 游戏对象 | ✅ 已完成 | `export_objects_mysql.py` |
| **Locks**     | 锁/分类  | ✅ 已完成 | `export_objects_mysql.py` |
| **Quests**    | 任务     | ✅ 已完成 | `export_quests.py`        |
| **Creatures** | NPC      | ✅ 已完成 | `export_creatures.py`     |
| **Factions**  | 阵营     | ✅ 已完成 | `export_factions.py`      |
| **Spells**    | 法术     | ✅ 已完成 | `export_spells.py`        |
| **Loot**      | 掉落表   | ✅ 已完成 | `export_loot.py`          |

### 待导入表

| 表名          | 描述     | 优先级 | 备注                    |
| ------------- | -------- | ------ | ----------------------- |
| `item_set`    | 套装     | 高     | 需要 `aowow_itemset`    |
| `spell_icons` | 法术图标 | 中     | 需要 `aowow_spellicons` |
| `zones`       | 区域地图 | 中     | 需要 `aowow_zones`      |

---

## 功能开发进度

### 后端 API (Go)

- [x] **Items**: 获取详情、搜索、分类建议
- [x] **Loot**: AtlasLoot 层级浏览、掉落查询
- [x] **Objects**: 分类浏览 (基于 locks)、搜索
- [x] **Quests**: 详情 API、搜索 API
- [x] **Creatures**: 详情 API、搜索 API
- [x] **Factions**: 列表 API
- [x] **Spells**: 详情 API

### 前端页面 (React)

- [x] **Loot Browser**: 完整的 AtlasLoot 浏览界面
- [x] **Objects Browser**: 分类浏览、详情展示
- [x] **Quest Browser**: 任务列表、详情页
- [x] **Creature Browser**: NPC 列表、详情、掉落查看
- [x] **Faction Browser**: 阵营列表
- [x] **Spell Browser**: 法术查询、详情页、描述解析 ($s1, $d)

## 下一步计划

1.  **后端**: 补充 Item Sets 的导入和 API。
2.  **前端**: 优化详情页的 UI 细节（如图标、掉落显示）。
3.  **系统**: 进一步优化 ETL 流程，支持更多数据表的增量更新。
