# 🎬 Simple Video Processor with FFmpeg + AWS (Free Tier Friendly)

Uma aplicação simples de exemplo que utiliza o **FFmpeg** para processar vídeos brutos (_raw_) e convertê-los para um formato compatível com **reprodutores de streaming de vídeo** modernos (ex: HLS ou MP4 otimizados).

O foco é demonstrar uma pipeline básica de processamento de vídeo com **recursos da AWS**, visando **hospedagem no Free Tier** de forma econômica e funcional.

---

## ☁️ Recursos Utilizados (AWS)

Este projeto faz uso dos seguintes serviços da Amazon Web Services:

- **Amazon EC2** – Para executar o processamento com FFmpeg
- **Amazon S3** – Armazenamento dos arquivos de entrada e saída
- **Amazon CloudFront** – Distribuição dos vídeos processados via CDN
- **Amazon SNS** – Notificações de eventos no processamento
- **Amazon DynamoDB** – Registro e controle do status do processamento

---

## 🔄 Fluxo Básico da Aplicação

```text
[Usuário envia vídeo bruto] 
        ↓
[Vídeo salvo no S3]
        ↓
[SNS notifica EC2]
        ↓
[EC2 executa FFmpeg → converte vídeo]
        ↓
[Vídeo convertido armazenado no S3]
        ↓
[CloudFront entrega vídeo processado via URL pública]
        ↓
[DynamoDB registra status e metadados do processamento]
