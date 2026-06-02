export type UsuarioSessao = {
  id: string;
  nome: string;
  documento: string;
  email: string;
  estrangeiro: boolean;
};

export function decodificarJWT(
  token: string
): Record<string, unknown> | null {
  try {
    const parts = token.split(".");
    if (parts.length !== 3) return null;
    const payload = parts[1];
    const base64 = payload.replace(/-/g, "+").replace(/_/g, "/");
    const jsonStr = decodeURIComponent(
      atob(base64)
        .split("")
        .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
        .join("")
    );
    return JSON.parse(jsonStr);
  } catch {
    return null;
  }
}

export function claimsParaUsuario(
  claims: Record<string, unknown>
): UsuarioSessao {
  return {
    id: (claims.sub as string) || "",
    nome: (claims.name as string) || (claims.real_name as string) || "",
    documento: (claims.name as string) || "",
    email: (claims.email as string) || "",
    estrangeiro: false,
  };
}
