package telegram

var (
	startMsg_RU = `
	🤖Привет!🤖
	
	Этот бот может помочь вам:
	
	1. *Сохранять ссылки*:
	   Используйте команду: */save <ссылка> <описание>*
	   Пример: */save https://example.com полезная ссылка*
	   _P. S. Чтобы получить подробную справку по работе этой команды, просто нажмите на нее 👉 /save 👈_
	
	2. *Получать сохранённые ссылки* по описанию:
	   Используйте команду: */get <описание>*
	   Пример: */get полезная ссылка*
	   Бот вернёт вам список ссылок, подходящих по описанию (если есть совпадения).
	   При выборе одной из ссылок, бот вернет вам сгенерированную новую.
	   Пример: *https://serverhost.com/gen/yourusername/generatedlink*
	   _P. S. Чтобы получить подробную справку по работе этой команды, просто нажмите на нее 👉 /get 👈_
	
	3. *Удалять сохранённые ссылки*:
	   Используйте команду: /delete <описание>
	   Пример: */delete полезная ссылка*
	   _P. S. Чтобы получить подробную справку по работе этой команды, просто нажмите на нее 👉 /delete 👈_
	
	4. *Общая справка*:
	   Для получения общей справки, нажмите на 👉 /help 👈
	
	Надеюсь, это поможет вам эффективно использовать моего бота!
	`

	saveMsg_RU = `
	*Итак, как сохранять ссылки? 🙌*
	Вам необходимо написать такую команду: */save <ссылка> <описание>*

	Примеры ссылок: 
	https://translate.yandex.ru/
	https://habr.com/ru/articles/
	https://github.com/

	В общем, любая ссылка, которая ведет на какую-то страницу.
	Через пробел после ссылки дайте *название для ссылки*. 
	Например, если вы хотите сохранить статью про кенгуру, то так и назовите ее: кенгуру.
	Запомнить вам надо, так как в дальнейшем поиск выполняется именно по вашему названию.
	`

	getMsg_RU = `
	*Как получить сохраненную ссылку?*
	Вам необходимо написать такую команду: */get <описание>*

	Пример:
	*/save https://en.wikipedia.org/wiki/Kangaroo kangaroo wiki info*
	...
	*/get kangaroo*
	ИЛИ
	*/get kang*
	ИЛИ
	*/get wiki*
	ИЛИ
	*/get info*
	...
	И так далее.
	`
)

var (
	startMsg_EN = `
	🤖Hello!🤖
	
	This bot can help you with the following:

	1. *Save links*:
	   Use the command: */save <link> <description>*
	   Example: */save https://example.com useful link*
	   _P.S. To get detailed instructions on how this command works, just click on it 👉 /save 👈_

	2. *Retrieve saved links* by description:
	   Use the command: */get <description>*
	   Example: */get useful link*
	   The bot will return a list of links that match the description (if any matches are found).
	   When you select one of the links, the bot will return a generated new one.
	   Example: *https://serverhost.com/gen/yourusername/generatedlink*
	   _P.S. To get detailed instructions on how this command works, just click on it 👉 /get 👈_

	3. *Delete saved links*:
	   Use the command: /delete <description>
	   Example: */delete useful link*
	   _P.S. To get detailed instructions on how this command works, just click on it 👉 /delete 👈_

	4. *General help*:
	   To get general help, click on 👉 /help 👈
	
	I hope this helps you make the most of my bot!
	`

	saveMsg_EN = `
	*So, how do you save links? 🙌*
	You need to type the following command: */save <link> <description>*

	Examples of links: 
	https://translate.yandex.com/
	https://habr.com/en/articles/
	https://github.com/

	In general, any link that leads to a specific page.
	After the link, give a *name for the link* separated by a space.
	For example, if you want to save an article about kangaroos, name it accordingly: *kangaroo*.
	Remember this name, as future searches will be performed using the name you provided.

	Example:
	*/save https://en.wikipedia.org/wiki/Kangaroo kangaroo wiki info*
	`

	getMsg_EN = `
	*How to get saved link?*
	You need to type the following command: */get <description>*

	Example:
	*/save https://en.wikipedia.org/wiki/Kangaroo kangaroo wiki info*
	...
	*/get kangaroo*
	OR
	*/get kang*
	OR
	*/get wiki*
	OR
	*/get info*
	...
	etc
	`
)
