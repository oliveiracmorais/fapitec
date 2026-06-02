import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const PUBLIC_PATHS = [
  "/",
  "/auth/callback",
  "/registrar",
  "/recuperar-senha",
  "/redefinir-senha",
];

const PROTECTED_PREFIXES = ["/dashboard", "/editais", "/perfil"];

export function middleware(request: NextRequest) {
  const token = request.cookies.get("fapitec_token")?.value;
  const { pathname } = request.nextUrl;

  const isPublic = PUBLIC_PATHS.some(
    (p) => pathname === p || pathname.startsWith(p + "?")
  );

  const isProtected = PROTECTED_PREFIXES.some((p) =>
    pathname.startsWith(p)
  );

  if (isProtected && !token) {
    const loginUrl = new URL("/", request.url);
    loginUrl.searchParams.set("redirect", pathname);
    return NextResponse.redirect(loginUrl);
  }

  if (pathname.startsWith("/api/auth/")) {
    return NextResponse.next();
  }

  if (pathname.startsWith("/api/v1/")) {
    const apiUrl =
      process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    const url = new URL(pathname, apiUrl);
    url.search = request.nextUrl.search;

    const requestHeaders = new Headers(request.headers);
    if (token) {
      requestHeaders.set("Authorization", `Bearer ${token}`);
    }

    return NextResponse.rewrite(url, {
      request: { headers: requestHeaders },
    });
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/((?!_next/static|_next/image|favicon.ico|logo-2.png|Logo-1.jpeg|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)",
  ],
};
