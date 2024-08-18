# tech_stack_processor.py

def extract_tech_stack_features(tech_stacks):
    """
    tech_stackの辞書リストからバイナリ特徴量を抽出します。

    Parameters:
    tech_stacks (list): 各辞書が技術スタックを表すリスト。

    Returns:
    list: 各技術スタックに対応するバイナリ特徴量を含む辞書のリスト。
    """
    
    # tech_stack辞書のすべてのユニークなキーと値のペアを特定
    unique_features = set()

    for stack in tech_stacks:
        for key, values in stack.items():
            for value in values:
                unique_features.add(f'{key}_{value}')

    # 集めたユニークな特徴量をリストに変換
    unique_features = list(unique_features)

    # 1つのtech_stackから特徴量を抽出する関数
    def _extract_features(tech_stack):
        features = {feature: 0 for feature in unique_features}
        
        for key, values in tech_stack.items():
            for value in values:
                feature_name = f'{key}_{value}'
                if feature_name in features:
                    features[feature_name] = 1
        
        return features

    # すべてのtech_stackを処理し、特徴量の辞書リストを返す
    return [_extract_features(stack) for stack in tech_stacks]
