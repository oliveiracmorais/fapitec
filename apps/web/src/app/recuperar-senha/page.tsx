"use client";

import { useState, FormEvent } from "react";
import Link from "next/link";
import Image from "next/image";

export default function RecuperarSenhaPage() {
  const [email, setEmail] = useState("");
  const [mensagem, setMensagem] = useState("");
  const [erro, setErro] = useState("");
  const [carregando, setCarregando] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setErro("");
    setMensagem("");
    setCarregando(true);

    try {
      const res = await fetch("/api/v1/solicitar-redefinicao-senha", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email }),
      });

      const data = await res.json();

      if (!res.ok) {
        setErro(data.erro || "Erro ao solicitar redefinição");
        return;
      }

      setMensagem(
        "Se o e-mail estiver cadastrado, você receberá um link de redefinição de senha."
      );
    } catch {
      setErro("Erro de conexão com o servidor");
    } finally {
      setCarregando(false);
    }
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-linear-to-br from-brand-900 via-brand-800 to-brand-700 px-4">
      <div className="w-full max-w-md rounded-2xl bg-white/95 p-8 shadow-2xl backdrop-blur">
        <div className="mb-8 flex flex-col items-center">
          <Image
            src="/logo-2.png"
            alt="FAPITEC-SE"
            width={211}
            height={60}
            className="mb-4"
            priority
          />
          <h1 className="text-xl font-bold text-brand-800">
            Recuperar Senha
          </h1>
          <p className="text-sm text-gray-600">
            Informe seu e-mail para receber um link de redefinição
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label
              htmlFor="email"
              className="flex items-center gap-1.5 text-sm font-medium text-gray-700"
            >
              <span>📧</span> E-mail
            </label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="seu@email.com"
              required
              className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
            />
          </div>

          {erro && (
            <div className="rounded-lg bg-red-50 p-3 text-sm text-red-700">
              {erro}
            </div>
          )}

          {mensagem && (
            <div className="rounded-lg bg-green-50 p-3 text-sm text-green-700">
              {mensagem}
            </div>
          )}

          <button
            type="submit"
            disabled={carregando}
            className="w-full rounded-lg bg-brand-600 px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-brand-700 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {carregando ? "Enviando..." : "Enviar link de redefinição"}
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-gray-600">
          <Link
            href="/"
            className="font-medium text-brand-600 hover:text-brand-700"
          >
            Voltar para o início
          </Link>
        </p>
      </div>
    </div>
  );
}
