# matrix
Stocks prices alerter

## dev
This project uses asdf for language versions management



🟢 Neo → Collector principal
Encargado de traer datos del mundo exterior: precios de acciones y noticias (equivalente a tus “Prices Collector” y “News Collector”).

🔵 Trinity → Ingestor / Processor
Se comunica por gRPC con Neo, recibe los datos crudos y los mete en el sistema (DB o Pub/Sub Bus). Puede encargarse de normalizar, enriquecer o limpiar los datos antes de guardarlos.

🟣 Morpheus → Alerter
Analiza la base de datos o el flujo de datos en el bus, detecta condiciones que requieren atención (por ejemplo, cambios bruscos en precios o noticias críticas), y emite un evento de alerta.

🟠 Tank → Notifier / Gateway
Toma los eventos generados por Morpheus y los distribuye a los canales externos (correo, Slack, Telegram, webhooks, etc.).