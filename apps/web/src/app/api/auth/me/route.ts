import { type NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  const token = request.cookies.get("fapitec_token")?.value;

  if (!token) {
    return NextResponse.json({ authenticated: false }, { status: 401 });
  }

  try {
    const payload = JSON.parse(
      decodeURIComponent(
        atob(
          token
            .split(".")[1]
            .replace(/-/g, "+")
            .replace(/_/g, "/")
        )
          .split("")
          .map(
            (c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2)
          )
          .join("")
      )
    );

    return NextResponse.json({
      authenticated: true,
      usuario: {
        id: payload.sub || "",
        nome: payload.name || payload.real_name || "",
        documento: payload.name || "",
        email: payload.email || "",
        estrangeiro: false,
      },
    });
  } catch {
    return NextResponse.json({ authenticated: false }, { status: 401 });
  }
}
