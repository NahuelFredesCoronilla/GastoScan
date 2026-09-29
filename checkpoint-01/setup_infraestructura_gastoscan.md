# Setup de Infraestructura y Servicios Cloud \- GastoScan

Este documento describe la arquitectura, configuración de red y el listado completo de componentes y servicios cloud requeridos para desplegar **GastoScan** en **Amazon Web Services (AWS)**.

---

## 1\. Red y Conectividad (Networking)

* **AWS VPC (Virtual Private Cloud)**:  
  * Espacio de red virtual aislado y dedicado para la plataforma.  
  * Rango CIDR recomendado: `/16` (ejemplo: `10.0.0.0/16`).  
* **Subredes Públicas (Mínimo 2 Availability Zones)**:  
  * Alojan componentes con exposición a internet: Application Load Balancer (ALB) y NAT Gateways.  
* **Subredes Privadas de Aplicación (Mínimo 2 Availability Zones)**:  
  * Alojan los contenedores backend en **AWS ECS Fargate** (Core API y Scan API).  
  * Sin asignación de IPs públicas.  
* **Subredes Privadas de Datos (Mínimo 2 Availability Zones)**:  
  * Alojan la instancia de base de datos **Amazon RDS PostgreSQL**.  
  * Totalmente aisladas del tráfico entrante de internet.  
* **Internet Gateway (IGW)**:  
  * Proporciona acceso y enrutamiento hacia las subredes públicas para recibir el tráfico de la app móvil.  
* **NAT Gateways / VPC Endpoints (AWS PrivateLink)**:  
  * Permiten a los servicios en subredes privadas comunicarse con servicios AWS gestionados (S3, SQS, Textract, Bedrock, CloudWatch) de forma segura.  
* **Security Groups**:  
  * `sg-alb`: Permite tráfico entrante HTTPS (puerto 443\) desde cualquier origen (`0.0.0.0/0`).  
  * `sg-ecs-core`: Permite tráfico HTTP entrante únicamente desde `sg-alb`.  
  * `sg-ecs-scan`: Permite tráfico saliente hacia SQS, S3, Textract, Bedrock y RDS.  
  * `sg-rds`: Permite tráfico PostgreSQL (puerto 5432\) únicamente desde `sg-ecs-core` y `sg-ecs-scan`.

---

## 2\. Capa de Entrada y Seguridad de Acceso

* **Application Load Balancer (ALB)**:  
  * Punto único de entrada para la aplicación cliente.  
  * Terminación TLS/HTTPS con certificado SSL/TLS gestionado mediante **AWS Certificate Manager (ACM)**.  
  * Enrutamiento de peticiones hacia el Target Group de Core API.  
* **Amazon Cognito (User Pool & App Client)**:  
  * Directorio de usuarios para registro, inicio de sesión y recuperación de contraseña.  
  * Emisión y validación de tokens JWT (JSON Web Tokens) consumidos por la aplicación Flutter.

---

## 3\. Cómputo y Contenedores (Backend Layer)

* **Amazon ECR (Elastic Container Registry)**:  
  * Repositorio privado para almacenar y versionar las imágenes Docker de:  
    * `gastoscan-core-api`  
    * `gastoscan-scan-api`  
* **AWS ECS (Elastic Container Service) Cluster**:  
  * Orquestador de contenedores configurado con capacidad **AWS Fargate** (serverless).  
* **Servicios ECS Fargate**:  
  * **Core API (Go)**:  
    * Servicio HTTP que gestiona la lógica de negocio, endpoints CRUD de gastos, categorización manual y orquestación de subidas.  
    * Integrado con el Application Load Balancer.  
  * **Scan API (Go)**:  
    * Worker en segundo plano que consume eventos desde la cola SQS.  
    * Encargado del pipeline de procesamiento inteligente de comprobantes.

---

## 4\. Persistencia y Almacenamiento

* **Amazon RDS para PostgreSQL**:  
  * Instancia gestionada de base de datos relacional (Multi-AZ opcional para alta disponibilidad).  
  * Almacena tablas de usuarios, comprobantes, ítems, comercios, categorías e historiales.  
  * Cifrado en reposo habilitado con AWS KMS.  
* **Amazon S3 (Simple Storage Service)**:  
  * Bucket privado para almacenar archivos multimedia (PDF, JPG, PNG) de comprobantes y tickets.  
  * Bloqueo total de acceso público (*Block Public Access* activado).  
  * Cifrado del lado del servidor habilitado (SSE-S3 o SSE-KMS).  
  * Mecanismo de subida/descarga mediante Presigned URLs o IAM Roles específicos de ECS.

---

## 5\. Mensajería y Desacoplamiento Asíncrono

* **Amazon SQS (Standard Queue)**:  
  * Cola principal de procesamiento de comprobantes (`gastoscan-receipt-processing-queue`).  
  * Desacopla la recepción del archivo en Core API del procesamiento con IA en Scan API.  
  * Configuración de *Visibility Timeout* adaptada al tiempo de respuesta de OCR y LLMs.  
* **Amazon SQS Dead-Letter Queue (DLQ)**:  
  * Cola de mensajes fallidos (`gastoscan-receipt-dlq`) tras superar el número máximo de reintentos (*maxReceiveCount*).

---

## 6\. Servicios de Inteligencia Artificial (Managed AI)

* **Amazon Textract**:  
  * Servicio de OCR para extracción de texto estructurado, campos clave-valor y tablas de facturas/recibos.  
* **Amazon Bedrock**:  
  * Acceso a Modelos Fundacionales (LLMs).  
  * Utilizado por Scan API para:  
    * Limpieza y normalización de los textos extraídos por Textract.  
    * Clasificación inteligente de categorías y comercios.  
    * Estructuración final de datos para inserción en la base de datos.

---

## 7\. Observabilidad, Monitoreo y Configuración

* **Amazon CloudWatch Logs**:  
  * Centralización de logs de los contenedores Core API y Scan API via driver `awslogs`.  
* **Amazon CloudWatch Metrics & Alarms**:  
  * Métricas de uso de CPU y memoria en ECS Fargate.  
  * Monitoreo de latencia y código de respuesta en el ALB.  
  * Alarma de mensajes en espera en la cola SQS (`ApproximateNumberOfMessagesVisible`).  
* **AWS Secrets Manager / SSM Parameter Store**:  
  * Almacenamiento centralizado y cifrado de secretos (credenciales de base de datos, llaves de API, configuraciones sensibles).

---

## 8\. Pipeline de Integración y Despliegue Continuo (CI/CD)

* **GitHub Actions**:  
  * Ejecución de pruebas unitarias y linters en Go.  
  * Construcción de imágenes Docker multi-stage.  
  * Autenticación en AWS mediante OIDC / IAM Roles asumidos (sin credenciales hardcodeadas).  
  * Push de imágenes a Amazon ECR y actualización de definiciones de tareas (`task-definition.json`) en AWS ECS Fargate.