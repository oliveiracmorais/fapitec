"use client";

import { useState, FormEvent } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

type UsuarioSessao = {
  id: number;
  nome: string;
  documento: string;
  email: string;
  estrangeiro: boolean;
};

function salvarSessao(usuario: UsuarioSessao) {
  localStorage.setItem(
    "sessao",
    JSON.stringify({ usuario, timestamp: Date.now() })
  );
}

export default function LoginPage() {
  const router = useRouter();
  const [cpf, setCpf] = useState("");
  const [senha, setSenha] = useState("");
  const [erro, setErro] = useState("");
  const [carregando, setCarregando] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setErro("");
    setCarregando(true);

    try {
      const res = await fetch("/api/v1/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ cpf, senha }),
      });

      const data = await res.json();

      if (!res.ok) {
        setErro(data.erro || "Erro ao autenticar");
        return;
      }

      salvarSessao(data as UsuarioSessao);
      router.push("/dashboard");
    } catch {
      setErro("Erro de conexão com o servidor");
    } finally {
      setCarregando(false);
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center px-4">
      <div className="w-full max-w-md rounded-xl bg-white p-8 shadow-lg">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold text-gray-900">FAPITEC-SE</h1>
          <p className="mt-1 text-sm text-gray-500">
            Plataforma integrada de gestão institucional
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label
              htmlFor="cpf"
              className="block text-sm font-medium text-gray-700"
            >
              CPF ou Passaporte
            </label>
            <input
              id="cpf"
              type="text"
              value={cpf}
              onChange={(e) => setCpf(e.target.value)}
              placeholder="000.000.000-00"
              required
              className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm focus:border-blue-600 focus:outline-none focus:ring-1 focus:ring-blue-600"
            />
          </div>

          <div>
            <label
              htmlFor="senha"
              className="block text-sm font-medium text-gray-700"
            >
              Senha
            </label>
            <input
              id="senha"
              type="password"
              value={senha}
              onChange={(e) => setSenha(e.target.value)}
              placeholder="Sua senha"
              required
              className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm focus:border-blue-600 focus:outline-none focus:ring-1 focus:ring-blue-600"
            />
          </div>

          {erro && (
            <div className="rounded-lg bg-red-50 p-3 text-sm text-red-700">
              {erro}
            </div>
          )}

          <button
            type="submit"
            disabled={carregando}
            className="w-full rounded-lg bg-blue-700 px-4 py-2.5 text-sm font-medium text-white hover:bg-blue-800 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {carregando ? "Entrando..." : "Entrar"}
          </button>
        </form>

        <div className="mt-6 space-y-3 text-center text-sm">
          <Link
            href="/recuperar-senha"
            className="block text-blue-700 hover:text-blue-800"
          >
            Esqueci minha senha
          </Link>
          <p className="text-gray-500">
            Não tem cadastro?{" "}
            <Link
              href="/cadastro"
              className="font-medium text-blue-700 hover:text-blue-800"
            >
              Clique aqui
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
