

University: [ITMO University](https://itmo.ru/ru/)
Faculty: [FICT](https://fict.itmo.ru)
Course: [Network programming](https://github.com/itmo-ict-faculty/introduction-in-routing)
Year: 2024/2025
Group: K3323
Author: Abdrakhmanov Martin Maratovich
Lab: Lab2
Date of create: 1.04.2025
Date of finished: - 


# Поднимем Netbox

![image](https://github.com/user-attachments/assets/57328434-c85d-4765-87c6-2b9cd2a1f5a3)


# Заполним данные о CHR

![image](https://github.com/user-attachments/assets/6231bdf5-7fbd-4ff1-ac86-4ee9186b5458)

К сожалению в бесплатной версии open vpn максимум 2 активных коннекта, так что ради выполнения данной лабораторной работы изменим ip на локальные, чтобы в дальнейшем подключится к микротикам.


![image](https://github.com/user-attachments/assets/414fd5b2-525b-4786-a37c-ff7390589c80)

Заранее добавим кастомные параметры, потом пригодятся.

![image](https://github.com/user-attachments/assets/ba08f36c-3484-4c5d-a929-80ec0c9f0876)

# Выведем все данные из Netbox 
![image](https://github.com/user-attachments/assets/cb959315-6f21-4368-8fee-2405555c0cc8)

Подробнее можно посмотреть в файле netbox_data.json

# Напишем сценарий настройки 2 CHR

![image](https://github.com/user-attachments/assets/d479d393-54e4-4f9d-8318-b185baf65f91)

Информацию он берет из кастомных полей, заполенных ранее.

Микротики переименовались (изменение ip адреса такое же, в рамках тестирования не изменялось)

![image](https://github.com/user-attachments/assets/ba646a1d-64b0-4065-8770-8d06d3f56818)



# Сценарий сбора серийных номеров

![image](https://github.com/user-attachments/assets/11ee5de1-cbc3-42d5-b5af-c45ba9d035e3)

Так как прошлым сценарием мы обновили имена устройств так, чтобы они соотвествовали соответсвенным в Netbox, то серийный номер можно вносить просто по именам устройств.

![image](https://github.com/user-attachments/assets/e1811167-8a37-4a34-a174-e21cddb9aaa4)

Проверяем

![image](https://github.com/user-attachments/assets/38cc0e48-c3a6-4d69-9a94-395c680d710f)






