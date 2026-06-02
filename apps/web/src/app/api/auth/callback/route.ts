import { type NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  try {
    const { code, state } = await request.json();

    if (!code) {
      return NextResponse.json({ erro: "code ausente" }, { status: 400 });
    }

    const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    const res = await fetch(`${apiUrl}/api/v1/auth/callback`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ code, state: state || "fapitec-state" }),
    });

    if (!res.ok) {
      const err = await res
        .json()
        .catch(() => ({ erro: "falha ao obter token" }));
      return NextResponse.json(err, { status: res.status });
    }

    const data = await res.json();
    const token = data.access_token;

    if (!token) {
      return NextResponse.json({ erro: "token nao recebido" }, { status: 502 });
    }

    const response = NextResponse.json({ success: true });
    response.cookies.set("fapitec_token", token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
      maxAge: 60 * 60 * 24,
    });

    return response;
  } catch {
    return NextResponse.json({ erro: "erro interno" }, { status: 500 });
  }
}
