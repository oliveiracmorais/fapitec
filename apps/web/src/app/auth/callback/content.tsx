"use client";

import { useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";

export default function CallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const code = searchParams.get("code");
    const state = searchParams.get("state");

    if (!code) {
      router.push("/?erro=code_ausente");
      return;
    }

    fetch("/api/auth/callback", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ code, state }),
    })
      .then((res) => {
        if (res.ok) router.push("/dashboard");
        else
          res
            .json()
            .then((d) =>
              router.push(
                `/?erro=${encodeURIComponent(d.erro || "falha_autenticacao")}`
              )
            );
      })
      .catch(() => router.push("/?erro=erro_conexao"));
  }, [router, searchParams]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-brand-900 via-brand-800 to-brand-700">
      <div className="rounded-xl bg-white p-8 shadow-2xl">
        <p className="text-gray-700">Autenticando...</p>
      </div>
    </div>
  );
}
