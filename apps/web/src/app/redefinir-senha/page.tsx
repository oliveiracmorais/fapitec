"use client";

import { Suspense, useState, FormEvent, useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import Image from "next/image";

function RedefinirSenhaForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const tokenFromUrl = searchParams.get("token") || "";

  const [token, setToken] = useState(tokenFromUrl);
  const [senha, setSenha] = useState("");
  const [confirmacaoSenha, setConfirmacaoSenha] = useState("");
  const [mensagem, setMensagem] = useState("");
  const [erro, setErro] = useState("");
  const [carregando, setCarregando] = useState(false);

  useEffect(() => {
    if (tokenFromUrl) {
      setToken(tokenFromUrl);
    }
  }, [tokenFromUrl]);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setErro("");
    setMensagem("");
    setCarregando(true);

    try {
      const res = await fetch("/api/v1/redefinir-senha", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ token, senha, confirmacao_senha: confirmacaoSenha }),
      });

      const data = await res.json();

      if (!res.ok) {
        setErro(data.erro || "Erro ao redefinir senha");
        return;
      }

      setMensagem("Senha redefinida com sucesso! Redirecionando para o login...");
      setTimeout(() => router.push("/"), 2000);
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
            Redefinir Senha
          </h1>
          <p className="text-sm text-gray-600">
            Digite sua nova senha
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label
              htmlFor="token"
              className="flex items-center gap-1.5 text-sm font-medium text-gray-700"
            >
              <span>🔑</span> Token de redefinição
            </label>
            <input
              id="token"
              type="text"
              value={token}
              onChange={(e) => setToken(e.target.value)}
              placeholder="Cole o token recebido por email"
              required
              className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
            />
          </div>

          <div>
            <label
              htmlFor="senha"
              className="flex items-center gap-1.5 text-sm font-medium text-gray-700"
            >
              <span>🔒</span> Nova senha
            </label>
            <input
              id="senha"
              type="password"
              value={senha}
              onChange={(e) => setSenha(e.target.value)}
              placeholder="Mínimo 8 caracteres, maiúscula, minúscula, número e especial"
              required
              minLength={8}
              className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm outline-none transition-colors focus:border-brand-600 focus:ring-1 focus:ring-brand-600"
            />
          </div>

          <div>
            <label
              htmlFor="confirmacaoSenha"
              className="flex items-center gap-1.5 text-sm font-medium text-gray-700"
            >
              <span>🔒</span> Confirmar nova senha
            </label>
            <input
              id="confirmacaoSenha"
              type="password"
              value={confirmacaoSenha}
              onChange={(e) => setConfirmacaoSenha(e.target.value)}
              placeholder="Repita a nova senha"
              required
              minLength={8}
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
            {carregando ? "Redefinindo..." : "Redefinir senha"}
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

export default function RedefinirSenhaPage() {
  return (
    <Suspense fallback={
      <div className="flex min-h-screen items-center justify-center bg-linear-to-br from-brand-900 via-brand-800 to-brand-700">
        <p className="text-white">Carregando...</p>
      </div>
    }>
      <RedefinirSenhaForm />
    </Suspense>
  );
}
