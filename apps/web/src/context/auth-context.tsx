"use client";

import {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
  type ReactNode,
} from "react";
import { useRouter } from "next/navigation";
import type { UsuarioSessao } from "../lib/auth";

type AuthState = {
  usuario: UsuarioSessao | null;
  carregando: boolean;
  loginCasdoor: () => void;
  logout: () => Promise<void>;
  refreshAuth: () => Promise<void>;
  definirUsuario: (u: UsuarioSessao) => void;
};

const AuthContext = createContext<AuthState>({
  usuario: null,
  carregando: true,
  loginCasdoor: () => {},
  logout: async () => {},
  refreshAuth: async () => {},
  definirUsuario: () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [usuario, setUsuario] = useState<UsuarioSessao | null>(null);
  const [carregando, setCarregando] = useState(true);
  const router = useRouter();

  useEffect(() => {
    const stored = localStorage.getItem("sessao");
    if (stored) {
      try {
        const parsed = JSON.parse(stored);
        const idade = Date.now() - parsed.timestamp;
        if (idade < 86400000 && parsed.usuario) {
          setUsuario(parsed.usuario as UsuarioSessao);
          setCarregando(false);
          return;
        }
      } catch {
        localStorage.removeItem("sessao");
      }
    }

    fetch("/api/auth/me")
      .then((res) => {
        if (!res.ok) throw new Error("nao autenticado");
        return res.json();
      })
      .then((data) => {
        if (data.authenticated) {
          setUsuario(data.usuario);
        }
      })
      .catch(() => setUsuario(null))
      .finally(() => setCarregando(false));
  }, []);

  const definirUsuario = useCallback((u: UsuarioSessao) => {
    setUsuario(u);
  }, []);

  const refreshAuth = useCallback(async () => {
    setCarregando(true);
    try {
      const res = await fetch("/api/auth/me");
      if (res.ok) {
        const data = await res.json();
        if (data.authenticated) {
          setUsuario(data.usuario);
          return;
        }
      }
      setUsuario(null);
    } catch {
      setUsuario(null);
    } finally {
      setCarregando(false);
    }
  }, []);

  const loginCasdoor = useCallback(() => {
    window.location.href = "/api/v1/auth/login";
  }, []);

  const logout = useCallback(async () => {
    await fetch("/api/auth/logout", { method: "POST" }).catch(() => {});
    localStorage.removeItem("sessao");
    setUsuario(null);
    router.push("/");
  }, [router]);

  return (
    <AuthContext.Provider
      value={{ usuario, carregando, loginCasdoor, logout, refreshAuth, definirUsuario }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
