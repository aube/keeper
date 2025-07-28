# gophkeeper-client

## сборка

make build или ./cmd/client/build.sh

Соберутся бинарники под Windows, Linux, Mac на CPU x86/ARM:

- keeper_darwin_amd64
- keeper_darwin_arm64
- keeper_linux_amd64
- keeper_linux_arm64
- keeper_windows_amd64
- keeper_windows_arm64


## команды запуска

**Прежде чем логиниться надо зарегистрироваться**

**Прежде чем что-то делать надо залогиниться**

### Регистрация
`keeper_linux_amd64 register -u username -p password -e email`

### Аутентификация
`keeper_linux_amd64 login -u username -p password`

### Шифрование
`keeper_linux_amd64 encrypt -u username -p password -i filepath -o filename`

Получает файл, шифрует, отправляет на сервер

### Шифрование банковской карты (любых текстовых данных)
`keeper_linux_amd64 card -u username -p password -n number -d date -v cvv`

Получает данные карты, шифрует, отправляет на сервер

### Дешифрование
`keeper_linux_amd64 decrypt -u username -p password -i filename -o filepath`

Получает файл с сервера, дешифрует, сохраняет по указанному пути

### Чтение банковской карты (любых текстовых данных)
`keeper_linux_amd64 readcard -u username -p password -n number`

Получает данные с сервера, дешифрует, выводит на экран

### Синхронизация данных с сервером
`keeper_linux_amd64 sync -u username`



# TUI (пользовательский интерфейс)
Запуск без команды
`keeper_linux_amd64`

