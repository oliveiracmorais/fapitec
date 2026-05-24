"use client";

import Link from "next/link";

export default function RecuperarSenhaPage() {
  return (
    <main className="shell">
      <header className="topbar">
        <div className="brand">
          <strong>FAPITEC-SE</strong>
          <span>Recuperação de senha</span>
        </div>
      </header>

      <section className="content">
        <div className="placeholder" style={{ maxWidth: 500, margin: "40px auto" }}>
          <h1>Recuperar senha</h1>
          <p>
            Informe seu e-mail cadastrado para receber um link de redefinição de senha.
          </p>

          <form
            onSubmit={async (e) => {
              e.preventDefault();
              const form = e.target as HTMLFormElement;
              const data = new FormData(form);
              const email = data.get("email") as string;

              try {
                const res = await fetch("/api/v1/solicitar-redefinicao-senha", {
                  method: "POST",
                  headers: { "Content-Type": "application/json" },
                  body: JSON.stringify({ email }),
                });
                const json = await res.json();
                alert(json.mensagem || json.erro);
              } catch {
                alert("Erro ao conectar com o servidor.");
              }
            }}
            style={{ display: "grid", gap: 12 }}
          >
            <input
              type="email"
              name="email"
              placeholder="Seu e-mail cadastrado"
              required
              style={{
                padding: "10px 14px",
                border: "1px solid #d8e0e7",
                borderRadius: 6,
                fontSize: 14,
              }}
            />
            <button
              type="submit"
              style={{
                padding: "10px 14px",
                border: "none",
                borderRadius: 6,
                background: "var(--fapitec-azul)",
                color: "#fff",
                fontSize: 14,
                fontWeight: 700,
                cursor: "pointer",
              }}
            >
              Enviar link de redefinição
            </button>
          </form>

          <p style={{ fontSize: 13, color: "var(--fapitec-cinza-700)" }}>
            <Link href="/" style={{ color: "var(--fapitec-azul)", fontWeight: 700 }}>
              Voltar para o início
            </Link>
          </p>
        </div>
      </section>
    </main>
  );
}
