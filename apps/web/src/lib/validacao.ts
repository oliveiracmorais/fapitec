export type ErrosFormulario = Record<string, string>;

export function validarNome(valor: string): string {
  if (!valor.trim()) return "Nome completo é obrigatório";
  if (valor.trim().length < 3) return "Nome deve ter no mínimo 3 caracteres";
  if (!/^[a-zA-ZÀ-ÿ\s]+$/.test(valor))
    return "Nome deve conter apenas letras e espaços";
  return "";
}

export function validarCPF(valor: string): string {
  const digitos = valor.replace(/\D/g, "");
  if (digitos.length !== 11) return "CPF deve ter 11 dígitos";
  if (/^(\d)\1{10}$/.test(digitos))
    return "CPF inválido: todos os dígitos são iguais";
  if (!validarDigitosCPF(digitos))
    return "CPF inválido: dígitos verificadores não conferem";
  return "";
}

function validarDigitosCPF(cpf: string): boolean {
  const d1 = calcularDigitoCPF(cpf.slice(0, 9), 10);
  const d2 = calcularDigitoCPF(cpf.slice(0, 9) + d1, 11);
  return d1 === Number(cpf[9]) && d2 === Number(cpf[10]);
}

function calcularDigitoCPF(base: string, pesoInicial: number): number {
  let soma = 0;
  for (let i = 0; i < base.length; i++) {
    soma += Number(base[i]) * (pesoInicial - i);
  }
  const resto = soma % 11;
  return resto < 2 ? 0 : 11 - resto;
}

export function validarPassaporte(valor: string): string {
  if (!valor.trim()) return "Passaporte é obrigatório";
  if (!/^[a-zA-Z0-9]{3,30}$/.test(valor))
    return "Passaporte deve ter entre 3 e 30 caracteres alfanuméricos";
  return "";
}

export function validarEmail(valor: string): string {
  if (!valor.trim()) return "E-mail é obrigatório";
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(valor)) return "E-mail inválido";
  return "";
}

export function validarConfirmacaoEmail(
  email: string,
  confirmacao: string
): string {
  if (!confirmacao.trim()) return "Confirmação de e-mail é obrigatória";
  if (email !== confirmacao)
    return "O e-mail deve ser IGUAL ao e-mail principal.";
  return "";
}

export function validarSenha(valor: string): string {
  if (!valor) return "Senha é obrigatória";
  if (valor.length < 8) return "Senha deve ter no mínimo 8 caracteres";
  if (!/[A-Z]/.test(valor))
    return "Senha deve conter pelo menos uma letra maiúscula";
  if (!/[a-z]/.test(valor))
    return "Senha deve conter pelo menos uma letra minúscula";
  if (!/[0-9]/.test(valor))
    return "Senha deve conter pelo menos um número";
  if (!/[^a-zA-Z0-9]/.test(valor))
    return "Senha deve conter pelo menos um caractere especial";
  return "";
}

export function validarConfirmacaoSenha(
  senha: string,
  confirmacao: string
): string {
  if (!confirmacao) return "Confirmação de senha é obrigatória";
  if (senha !== confirmacao)
    return "A senha deve ser IGUAL à primeira senha fornecida.";
  return "";
}

export function formatarCPF(valor: string): string {
  const digitos = valor.replace(/\D/g, "").slice(0, 11);
  if (digitos.length <= 3) return digitos;
  if (digitos.length <= 6)
    return `${digitos.slice(0, 3)}.${digitos.slice(3)}`;
  if (digitos.length <= 9)
    return `${digitos.slice(0, 3)}.${digitos.slice(3, 6)}.${digitos.slice(6)}`;
  return `${digitos.slice(0, 3)}.${digitos.slice(3, 6)}.${digitos.slice(6, 9)}-${digitos.slice(9)}`;
}
