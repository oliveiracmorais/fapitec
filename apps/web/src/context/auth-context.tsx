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
};

const AuthContext = createContext<AuthState>({
  usuario: null,
  carregando: true,
  loginCasdoor: () => {},
  logout: async () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [usuario, setUsuario] = useState<UsuarioSessao | null>(null);
  const [carregando, setCarregando] = useState(true);
  const router = useRouter();

  useEffect(() => {
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

  const loginCasdoor = useCallback(() => {
    window.location.href = "/api/v1/auth/login";
  }, []);

  const logout = useCallback(async () => {
    await fetch("/api/auth/logout", { method: "POST" });
    setUsuario(null);
    router.push("/");
  }, [router]);

  return (
    <AuthContext.Provider
      value={{ usuario, carregando, loginCasdoor, logout }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
