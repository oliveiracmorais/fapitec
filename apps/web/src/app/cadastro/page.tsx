"use client";

import { useState, FormEvent } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

export default function CadastroPage() {
  const router = useRouter();
  const [form, setForm] = useState({
    nome: "",
    email: "",
    confirmacaoEmail: "",
    senha: "",
    confirmacaoSenha: "",
    estrangeiro: false,
    documento: "",
  });
  const [erro, setErro] = useState("");
  const [carregando, setCarregando] = useState(false);

  function atualizar(chave: string, valor: string | boolean) {
    setForm((prev) => ({ ...prev, [chave]: valor }));
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setErro("");

    if (form.email !== form.confirmacaoEmail) {
      setErro("O e-mail deve ser IGUAL ao e-mail principal.");
      return;
    }
    if (form.senha !== form.confirmacaoSenha) {
      setErro("A senha deve ser IGUAL à primeira senha fornecida.");
      return;
    }

    setCarregando(true);

    try {
      const body: Record<string, unknown> = {
        nome: form.nome,
        email: form.email,
        confirmacao_email: form.confirmacaoEmail,
        senha: form.senha,
        confirmacao_senha: form.confirmacaoSenha,
        estrangeiro: form.estrangeiro,
      };

      if (form.estrangeiro) {
        body.passaporte = form.documento;
      } else {
        body.cpf = form.documento;
      }

      const res = await fetch("/api/v1/cadastro", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });

      const data = await res.json();

      if (!res.ok) {
        setErro(data.erro || "Erro ao cadastrar");
        return;
      }

      router.push("/?cadastro=ok");
    } catch {
      setErro("Erro de conexão com o servidor");
    } finally {
      setCarregando(false);
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center px-4 py-8">
      <div className="w-full max-w-md rounded-xl bg-white p-8 shadow-lg">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold text-gray-900">Criar Conta</h1>
          <p className="mt-1 text-sm text-gray-500">
            Preencha os dados para se cadastrar
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label
              htmlFor="nome"
              className="block text-sm font-medium text-gray-700"
            >
              Nome Completo
            </label>
            <input
              id="nome"
              type="text"
              value={form.nome}
              onChange={(e) => atualizar("nome", e.target.value)}
              placeholder="Seu nome completo"
              required
              className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm focus:border-blue-600 focus:outline-none focus:ring-1 focus:ring-blue-600"
            />
          </div>

          <div className="flex items-center gap-2">
            <input
              id="estrangeiro"
              type="checkbox"
              checked={form.estrangeiro}
              onChange={(e) => atualizar("estrangeiro", e.target.checked)}
              className="h-4 w-4 rounded border-gray-300 text-blue-700"
            />
            <label htmlFor="estrangeiro" className="text-sm text-gray-700">
              Sou estrangeiro
            </label>
          </div>

          <div>
            <label
              htmlFor="documento"
              className="block text-sm font-medium text-gray-700"
            >
              {form.estrangeiro ? "Passaporte" : "CPF"}
            </label>
            <input
              id="documento"
              type="text"
              value={form.documento}
              onChange={(e) => atualizar("documento", e.target.value)}
              placeholder={
                form.estrangeiro ? "Número do passaporte" : "000.000.000-00"
              }
              required
              className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm focus:border-blue-600 focus:outline-none focus:ring-1 focus:ring-blue-600"
            />
          </div>

          <div>
            <label
              htmlFor="email"
              className="block text-sm font-medium text-gray-700"
            >
              E-mail
            </label>
            <input
              id="email"
              type="email"
              value={form.email}
              onChange={(e) => atualizar("email", e.target.value)}
              placeholder="seu@email.com"
              required
              className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm focus:border-blue-600 focus:outline-none focus:ring-1 focus:ring-blue-600"
            />
          </div>

          <div>
            <label
              htmlFor="confirmacaoEmail"
              className="block text-sm font-medium text-gray-700"
            >
              Confirmação de E-mail
            </label>
            <input
              id="confirmacaoEmail"
              type="email"
              value={form.confirmacaoEmail}
              onChange={(e) => atualizar("confirmacaoEmail", e.target.value)}
              placeholder="Confirme seu e-mail"
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
              value={form.senha}
              onChange={(e) => atualizar("senha", e.target.value)}
              placeholder="Mínimo 8 caracteres"
              required
              className="mt-1 block w-full rounded-lg border border-gray-300 px-4 py-2.5 text-sm focus:border-blue-600 focus:outline-none focus:ring-1 focus:ring-blue-600"
            />
            <p className="mt-1 text-xs text-gray-400">
              Mínimo 8 caracteres, com letras, números e símbolos
            </p>
          </div>

          <div>
            <label
              htmlFor="confirmacaoSenha"
              className="block text-sm font-medium text-gray-700"
            >
              Confirmação de Senha
            </label>
            <input
              id="confirmacaoSenha"
              type="password"
              value={form.confirmacaoSenha}
              onChange={(e) => atualizar("confirmacaoSenha", e.target.value)}
              placeholder="Repita a senha"
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
            {carregando ? "Cadastrando..." : "Cadastrar"}
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-gray-500">
          Já tem conta?{" "}
          <Link
            href="/"
            className="font-medium text-blue-700 hover:text-blue-800"
          >
            Faça login
          </Link>
        </p>
      </div>
    </div>
  );
}
