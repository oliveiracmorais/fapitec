"use client";

import { useState, FormEvent, useEffect, useRef } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import Image from "next/image";
import CampoFormulario from "../components/campo-formulario";
import { useAuth } from "../context/auth-context";

declare global {
  interface Window {
    turnstile?: {
      render: (
        container: HTMLElement,
        opts: {
          sitekey: string;
          callback: (token: string) => void;
          "expired-callback"?: () => void;
          "error-callback"?: () => void;
        }
      ) => void;
    };
  }
}

const SITEKEY =
  process.env.NEXT_PUBLIC_TURNSTILE_SITEKEY ||
  "1x00000000000000000000000000000000AA";

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
  const { definirUsuario } = useAuth();
  const [cpf, setCpf] = useState("");
  const [senha, setSenha] = useState("");
  const [erro, setErro] = useState("");
  const [carregando, setCarregando] = useState(false);
  const [tentativas, setTentativas] = useState(0);
  const [captchaToken, setCaptchaToken] = useState("");
  const [captchaIndisponivel, setCaptchaIndisponivel] = useState(false);

  const captchaRef = useRef<HTMLDivElement>(null);

  const mostrarCaptcha = tentativas >= 3;

  useEffect(() => {
    const el = captchaRef.current;
    if (!mostrarCaptcha || !el || el.children.length > 0) return;

    function renderizar() {
      if (!window.turnstile) return;
      try {
        window.turnstile.render(el!, {
          sitekey: SITEKEY,
          callback: (token) => {
            setCaptchaIndisponivel(false);
            setCaptchaToken(token);
          },
          "expired-callback": () => setCaptchaToken(""),
          "error-callback": () => setCaptchaIndisponivel(true),
        });
      } catch {
        console.warn("[captcha] erro ao renderizar turnstile");
      }
    }

    function startPolling() {
      let tentativasPoll = 0;
      const timer = setInterval(() => {
        tentativasPoll++;
        if (window.turnstile) {
          clearInterval(timer);
          renderizar();
        } else if (tentativasPoll > 50) {
          clearInterval(timer);
          console.error("[captcha] script turnstile nao inicializou apos 10s");
          setCaptchaIndisponivel(true);
        }
      }, 200);
      return timer;
    }

    if (window.turnstile) {
      renderizar();
      return;
    }

    if (!document.querySelector('script[src*="turnstile"]')) {
      const s = document.createElement("script");
      s.src = "https://challenges.cloudflare.com/turnstile/v0/api.js";
      s.async = true;
      document.head.appendChild(s);
    }

    const timer = startPolling();
    return () => clearInterval(timer);
  }, [mostrarCaptcha]);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (mostrarCaptcha && !captchaToken && !captchaIndisponivel) return;
    setErro("");
    setCarregando(true);

    try {
      const body: Record<string, string> = { cpf, senha };
      if (mostrarCaptcha && !captchaIndisponivel) {
        body.captcha_token = captchaToken;
      }

      const res = await fetch("/api/v1/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });

      const data = await res.json();

      if (!res.ok) {
        if (data.erro?.includes("captcha")) {
          setTentativas(3);
        } else {
          setTentativas((prev) => prev + 1);
        }
        setErro(data.erro || "Erro ao autenticar");
        return;
      }

      setTentativas(0);
      setCaptchaToken("");
      definirUsuario(data as any);
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
          <p className="text-sm text-gray-600">
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

          {captchaIndisponivel && mostrarCaptcha && (
            <div className="rounded-lg bg-yellow-50 p-3 text-sm text-yellow-800">
              Captcha indisponível — prossiga sem verificação
            </div>
          )}

          <div
            ref={captchaRef}
            className={`flex justify-center transition-all duration-300 ${
              mostrarCaptcha
                ? "max-h-20 opacity-100"
                : "max-h-0 overflow-hidden opacity-0"
            }`}
          />

          <button
            type="submit"
            disabled={carregando || (mostrarCaptcha && !captchaToken && !captchaIndisponivel)}
            className="w-full rounded-lg bg-brand-600 px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-brand-700 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {carregando ? "Entrando..." : "Entrar"}
          </button>
        </form>

        <div className="mt-6 space-y-3 text-center text-sm">
          <div className="relative mb-4">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-gray-300" />
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="bg-white px-2 text-gray-500">ou</span>
            </div>
          </div>
          <button
            type="button"
            onClick={() => (window.location.href = "/api/v1/auth/login")}
            className="w-full rounded-lg border border-brand-300 bg-white px-4 py-2.5 text-sm font-medium text-brand-700 transition-colors hover:bg-brand-50"
          >
            Entrar com FAPITEC (GovBr)
          </button>
          <Link
            href="/recuperar-senha"
            className="block pt-4 text-brand-600 hover:text-brand-700"
          >
            Esqueci minha senha
          </Link>
          <p className="text-gray-600">
            Não tem cadastro?{" "}
            <Link
              href="/registrar"
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
