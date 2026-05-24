"use client";

import { useState, FormEvent } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import Image from "next/image";
import CampoFormulario from "../components/campo-formulario";

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
    <div className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-br from-brand-900 via-brand-800 to-brand-700 px-4">
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
          <h1 className="text-xl font-bold text-brand-800">FAPITEC-SE</h1>
          <p className="text-sm text-gray-500">
            Plataforma integrada de gestão institucional
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-5">
          <CampoFormulario
            id="cpf"
            label="CPF ou Passaporte"
            placeholder="000.000.000-00"
            value={cpf}
            onChange={(e) => setCpf(e.target.value)}
            icone="🆔"
            required
          />

          <CampoFormulario
            id="senha"
            label="Senha"
            placeholder="Sua senha"
            value={senha}
            onChange={(e) => setSenha(e.target.value)}
            icone="🔒"
            tipoSenha
            required
          />

          {erro && (
            <div className="rounded-lg bg-red-50 p-3 text-sm text-red-700">
              {erro}
            </div>
          )}

          <button
            type="submit"
            disabled={carregando}
            className="w-full rounded-lg bg-brand-600 px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-brand-700 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {carregando ? "Entrando..." : "Entrar"}
          </button>
        </form>

        <div className="mt-6 space-y-3 text-center text-sm">
          <Link
            href="/recuperar-senha"
            className="block text-brand-600 hover:text-brand-700"
          >
            Esqueci minha senha
          </Link>
          <p className="text-gray-500">
            Não tem cadastro?{" "}
            <Link
              href="/cadastro"
              className="font-medium text-brand-600 hover:text-brand-700"
            >
              Clique aqui
            </Link>
          </p>
        </div>
      </div>

      <p className="mt-6 text-center text-xs text-white/60">
        Fundação de Apoio à Pesquisa e à Inovação Tecnológica do Estado de
        Sergipe
      </p>
    </div>
  );
}
