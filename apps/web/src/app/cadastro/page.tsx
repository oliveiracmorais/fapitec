"use client";

import { useState, FormEvent, useCallback } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import CampoFormulario from "../../components/campo-formulario";
import {
  validarNome,
  validarCPF,
  validarPassaporte,
  validarEmail,
  validarConfirmacaoEmail,
  validarSenha,
  validarConfirmacaoSenha,
  formatarCPF,
  type ErrosFormulario,
} from "../../lib/validacao";

type FormState = {
  nome: string;
  email: string;
  confirmacaoEmail: string;
  senha: string;
  confirmacaoSenha: string;
  estrangeiro: boolean;
  documento: string;
};

type TouchedState = Record<string, boolean>;

function validarFormulario(form: FormState): ErrosFormulario {
  const erros: ErrosFormulario = {};
  erros.nome = validarNome(form.nome);
  if (form.estrangeiro) {
    erros.documento = validarPassaporte(form.documento);
  } else {
    erros.documento = validarCPF(form.documento);
  }
  erros.email = validarEmail(form.email);
  erros.confirmacaoEmail = validarConfirmacaoEmail(
    form.email,
    form.confirmacaoEmail
  );
  erros.senha = validarSenha(form.senha);
  erros.confirmacaoSenha = validarConfirmacaoSenha(
    form.senha,
    form.confirmacaoSenha
  );
  return erros;
}

export default function CadastroPage() {
  const router = useRouter();
  const [form, setForm] = useState<FormState>({
    nome: "",
    email: "",
    confirmacaoEmail: "",
    senha: "",
    confirmacaoSenha: "",
    estrangeiro: false,
    documento: "",
  });
  const [touched, setTouched] = useState<TouchedState>({});
  const [erroServidor, setErroServidor] = useState("");
  const [carregando, setCarregando] = useState(false);

  const erros = validarFormulario(form);
  const valido = Object.values(erros).every((e) => !e);

  function atualizar(chave: keyof FormState, valor: string | boolean) {
    setForm((prev) => ({ ...prev, [chave]: valor }));
  }

  function handleBlur(chave: string) {
    setTouched((prev) => ({ ...prev, [chave]: true }));
  }

  const handleSubmit = useCallback(
    async (e: FormEvent) => {
      e.preventDefault();
      setTouched({
        nome: true,
        documento: true,
        email: true,
        confirmacaoEmail: true,
        senha: true,
        confirmacaoSenha: true,
      });

      if (!valido) return;

      setErroServidor("");
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
          setErroServidor(data.erro || "Erro ao cadastrar");
          return;
        }

        router.push("/?cadastro=ok");
      } catch {
        setErroServidor("Erro de conexão com o servidor");
      } finally {
        setCarregando(false);
      }
    },
    [form, valido, router]
  );

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
          <CampoFormulario
            id="nome"
            label="Nome Completo"
            placeholder="Seu nome completo"
            value={form.nome}
            onChange={(e) => atualizar("nome", e.target.value)}
            onBlur={() => handleBlur("nome")}
            icone="👤"
            erro={erros.nome}
            touched={touched.nome}
            required
          />

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

          <CampoFormulario
            id="documento"
            label={form.estrangeiro ? "Passaporte" : "CPF"}
            placeholder={form.estrangeiro ? "Número do passaporte" : "000.000.000-00"}
            value={form.documento}
            onChange={(e) => {
              const valor = form.estrangeiro
                ? e.target.value
                : formatarCPF(e.target.value);
              atualizar("documento", valor);
            }}
            onBlur={() => handleBlur("documento")}
            icone={form.estrangeiro ? "🛂" : "🆔"}
            erro={erros.documento}
            touched={touched.documento}
            required
          />

          <CampoFormulario
            id="email"
            label="E-mail"
            type="email"
            placeholder="seu@email.com"
            value={form.email}
            onChange={(e) => atualizar("email", e.target.value)}
            onBlur={() => handleBlur("email")}
            icone="📧"
            erro={erros.email}
            touched={touched.email}
            required
          />

          <CampoFormulario
            id="confirmacaoEmail"
            label="Confirmação de E-mail"
            type="email"
            placeholder="Confirme seu e-mail"
            value={form.confirmacaoEmail}
            onChange={(e) => atualizar("confirmacaoEmail", e.target.value)}
            onBlur={() => handleBlur("confirmacaoEmail")}
            icone="📧"
            erro={erros.confirmacaoEmail}
            touched={touched.confirmacaoEmail}
            required
          />

          <CampoFormulario
            id="senha"
            label="Senha"
            placeholder="Mínimo 8 caracteres"
            value={form.senha}
            onChange={(e) => atualizar("senha", e.target.value)}
            onBlur={() => handleBlur("senha")}
            icone="🔒"
            tipoSenha
            erro={erros.senha}
            touched={touched.senha}
            required
          />
          <p className="-mt-2 text-xs text-gray-400">
            Mínimo 8 caracteres, com letras maiúsculas, minúsculas, números e
            símbolos
          </p>

          <CampoFormulario
            id="confirmacaoSenha"
            label="Confirmação de Senha"
            placeholder="Repita a senha"
            value={form.confirmacaoSenha}
            onChange={(e) => atualizar("confirmacaoSenha", e.target.value)}
            onBlur={() => handleBlur("confirmacaoSenha")}
            icone="🔒"
            tipoSenha
            erro={erros.confirmacaoSenha}
            touched={touched.confirmacaoSenha}
            required
          />

          {erroServidor && (
            <div className="rounded-lg bg-red-50 p-3 text-sm text-red-700">
              {erroServidor}
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
