#!/bin/bash

# 告知ブロック
announcement_block='    {{if (HasAnnouncement)}}\n    <div style="padding: 15px; margin: 0 0 20px 0; border-radius: 6px; background-color: {{(GetAnnouncement).BgColor}}; color: {{(GetAnnouncement).TextColor}}; font-size: 16px; font-weight: bold; text-align: center; border: 2px solid {{(GetAnnouncement).TextColor}}; word-wrap: break-word;">\n      {{(GetAnnouncement).Message}}\n    </div>\n    {{end}}\n'

# 処理済みテンプレート（スキップ対象）
skip_files=("top.gtpl" "currentevents.gtpl" "announcement-form.gtpl" "footer.gtpl" "param-local.gtpl" "param-global.gtpl" "param-event0.gtpl" "param-event1.gtpl" "param-event2.gtpl" "param-eventc.gtpl")

# 全.gtplファイルに対してループ
for file in *.gtpl; do
    # スキップ対象に含まれるか確認
    skip=0
    for skip_file in "${skip_files[@]}"; do
        if [ "$file" = "$skip_file" ]; then
            skip=1
            break
        fi
    done
    
    if [ $skip -eq 1 ]; then
        echo "Skipping $file"
        continue
    fi
    
    # <body> タグが存在するか確認
    if grep -q "<body>" "$file"; then
        # 既に告知ブロックが含まれているか確認
        if grep -q "HasAnnouncement" "$file"; then
            echo "Already has announcement block: $file"
            continue
        fi
        
        # <body> タグの直後に告知ブロックを挿入
        sed -i "/<body>/a\\$announcement_block" "$file"
        echo "Updated: $file"
    else
        echo "No <body> tag: $file"
    fi
done
