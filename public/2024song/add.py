import json

def add_to_json_file(file_path):
    with open(file_path, 'r', encoding='utf-8') as file:
        data = json.load(file)

    new_id = input("ID：")
    new_name = input("歌名：")
    new_singer = input("歌手：")
    new_path = "./music/"+new_name+".flac"

    new_object = {
        "id": new_id,
        "name": new_name,
        "singer": new_singer,
        "path": new_path
    }

    data.append(new_object)

    with open(file_path, 'w', encoding='utf-8') as file:
        json.dump(data, file, ensure_ascii=False, indent=4)
    print("Success")

file_path = 'playlist.json'
add_to_json_file(file_path)
