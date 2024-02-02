# Microsserviços em Go
Este projeto consiste em quatro microsserviços escritos em Go, cada um com uma responsabilidade específica. Eles são executados em um ambiente Docker usando o Docker Compose.

## Serviços

### 1. Auth Service

O Auth Service é responsável pela autenticação e gerenciamento de contas de usuário, utilizando tokens JWT.

Endpoints:
- `POST /user/login`: Realiza login e emite um token JWT.
- `POST /user`: Cria uma nova conta de usuário.

### 2. Logger Service
 
O Logger Service recebe logs dos outros serviços e os armazena no MongoDB.

### 3. Listener Service

O Listener Service funciona como um message bus, recebendo mensagens do RabbitMQ e encaminhando-as para os serviços correspondentes.

### 4. Mail Service

O Mail Service envia e-mails e é acionado quando o Auth Service envia uma mensagem para o Listener Service, que, por sua vez, direciona a mensagem para o Mail Service.

## Banco de Dados
O banco de dados utilizado para armazenar informações de login é o MySQL.

## Executando a Aplicação
Certifique-se de ter o Docker e o Docker Compose instalados em sua máquina. Em seguida, execute o seguinte comando na pasta /project do projeto:

```
docker-compose up --build -d
```

Isso iniciará todos os serviços e suas dependências.

## Configuração
Os serviços têm suas próprias configurações definidas em arquivos .env. Certifique-se de configurar corretamente esses arquivos para atender às suas necessidades.

## Tecnologias Utilizadas
- Go (Chi, outras bibliotecas)
- MySQL
- MongoDB
- RabbitMQ
- Docker / Docker Compose

## Contribuição
Contribuições são bem-vindas! Sinta-se à vontade para abrir problemas (issues) e enviar pull requests.

## Licença
Este projeto está licenciado sob a Licença MIT - consulte o arquivo LICENSE para obter detalhes.

Este README é um exemplo básico. Adapte-o conforme necessário, incluindo informações mais detalhadas, instruções específicas de configuração, dependências do sistema e qualquer outra coisa relevante para os usuários ou desenvolvedores que interagem com seus microsserviços.