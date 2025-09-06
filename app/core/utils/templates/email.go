package templates

import "fmt"

func CodeViaEmail(code string) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
			<html lang="pt-BR">
			<head>
				<meta charset="UTF-8" />
				<title>Seu código de autenticação - DecoreAgora</title>
				<meta name="viewport" content="width=device-width, initial-scale=1.0" />
			</head>
			<body
				style="margin:0;padding:0;background-color:#f9fafb;font-family:Arial,Helvetica,sans-serif;color:#1f2937;"
			>
				<table
				role="presentation"
				cellpadding="0"
				cellspacing="0"
				width="100%%"
				style="background-color:#f9fafb;padding:40px 0;"
				>
				<tr>
					<td align="center">
					<table
						role="presentation"
						cellpadding="0"
						cellspacing="0"
						width="100%%"
						style="max-width:600px;background-color:#ffffff;border-radius:16px;box-shadow:0 4px 12px rgba(0,0,0,0.05);overflow:hidden;"
					>
						<!-- Header -->
						<tr>
						<td
							style="background:linear-gradient(135deg,#9333ea,#ec4899);padding:24px;text-align:center;"
						>
							<h1
							style="margin:0;font-size:24px;font-weight:bold;color:#ffffff;"
							>
							DecoreAgora
							</h1>
						</td>
						</tr>

						<!-- Body -->
						<tr>
						<td style="padding:32px 24px;">
							<h2
							style="margin-top:0;font-size:20px;font-weight:700;color:#111827;text-align:center;"
							>
							Seu código de autenticação
							</h2>
							<p
							style="font-size:16px;line-height:1.5;color:#374151;text-align:center;"
							>
							Use o código abaixo para entrar com segurança na sua conta:
							</p>

							<!-- Code box -->
							<div
							style="margin:24px auto;text-align:center;padding:20px;background:linear-gradient(135deg,#8b5cf6,#ec4899);color:#ffffff;font-size:28px;font-weight:bold;letter-spacing:4px;border-radius:12px;max-width:300px;"
							>
							%s
							</div>

							<p
							style="font-size:14px;line-height:1.5;color:#6b7280;text-align:center;"
							>
							Este código expira em <strong>10 minutos</strong>.<br />
							Se você não solicitou este acesso, ignore este e-mail.
							</p>
						</td>
						</tr>

						<!-- Footer -->
						<tr>
						<td
							style="background-color:#f3f4f6;padding:20px;text-align:center;font-size:12px;color:#6b7280;"
						>
							© 2025 DecoreAgora. Todos os direitos reservados.
						</td>
						</tr>
					</table>
					</td>
				</tr>
				</table>
			</body>
			</html>
	`, code)
}
