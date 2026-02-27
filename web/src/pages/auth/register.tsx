import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { toast } from 'sonner'
import { authApi } from '@/services/auth'
import { RegisterRequest } from '@/types/auth'

const registerSchema = z.object({
  username: z.string().min(3, 'Username must be at least 3 characters'),
  email: z.string().email('Invalid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords don't match",
  path: ['confirmPassword'],
})

type RegisterFormData = z.infer<typeof registerSchema>

export function RegisterPage() {
  const navigate = useNavigate()
  const [isLoading, setIsLoading] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
  })

  const onSubmit = async (data: RegisterFormData) => {
    setIsLoading(true)
    try {
      const registerData: RegisterRequest = {
        username: data.username,
        email: data.email,
        password: data.password,
      }
      
      // Direct use of service; request interceptor handles code!==0 by throwing error
      await authApi.register(registerData)
      
      toast.success('Registration successful! Please sign in.')
      navigate('/login')
    } catch (error: any) {
      console.error('Registration error:', error)
      // Error is already handled/formatted by request interceptor but we can show extra toast here if needed
      // or let the global handler do it. For now, specific toast:
      const msg = error?.message || 'Registration failed';
      toast.error(msg);
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-[#efefef] px-4 py-8 sm:px-6">
      <div className="mx-auto flex min-h-[calc(100vh-4rem)] w-full max-w-md items-center justify-center">
        <form onSubmit={handleSubmit(onSubmit)} className="w-full space-y-4">
          <div>
            <h1 className="text-[2.1rem] font-semibold tracking-tight text-[#1c1d20]">
              Create your account
            </h1>
          </div>

          <div className="space-y-4">
            <div className="space-y-2.5">
              <Label htmlFor="username" className="text-[1.05rem] font-medium text-[#1f2022]">
                Username
              </Label>
              <Input
                id="username"
                type="text"
                autoComplete="username"
                placeholder="johndoe"
                {...register('username')}
                className="h-12 rounded-md border-[#d7d7d7] bg-[#ececec] px-4 text-base text-[#1f2022] placeholder:text-[#a0a0a0] focus-visible:ring-0"
              />
              {errors.username && (
                <p className="text-sm text-red-600">{errors.username.message}</p>
              )}
            </div>

            <div className="space-y-2.5">
              <Label htmlFor="email" className="text-[1.05rem] font-medium text-[#1f2022]">
                Email
              </Label>
              <Input
                id="email"
                type="email"
                autoComplete="email"
                placeholder="me@somewhere.com"
                {...register('email')}
                className="h-12 rounded-md border-[#d7d7d7] bg-[#ececec] px-4 text-base text-[#1f2022] placeholder:text-[#a0a0a0] focus-visible:ring-0"
              />
              {errors.email && (
                <p className="text-sm text-red-600">{errors.email.message}</p>
              )}
            </div>

            <div className="space-y-2.5">
              <Label htmlFor="password" className="text-[1.05rem] font-medium text-[#1f2022]">
                Password
              </Label>
              <Input
                id="password"
                type="password"
                autoComplete="new-password"
                placeholder="••••••••••"
                {...register('password')}
                className="h-12 rounded-md border-[#d7d7d7] bg-[#ececec] px-4 text-base text-[#1f2022] placeholder:text-[#8f8f8f] focus-visible:ring-0"
              />
              {errors.password && (
                <p className="text-sm text-red-600">{errors.password.message}</p>
              )}
            </div>

            <div className="space-y-2.5">
              <Label htmlFor="confirmPassword" className="text-[1.05rem] font-medium text-[#1f2022]">
                Confirm Password
              </Label>
              <Input
                id="confirmPassword"
                type="password"
                autoComplete="new-password"
                placeholder="••••••••••"
                {...register('confirmPassword')}
                className="h-12 rounded-md border-[#d7d7d7] bg-[#ececec] px-4 text-base text-[#1f2022] placeholder:text-[#8f8f8f] focus-visible:ring-0"
              />
              {errors.confirmPassword && (
                <p className="text-sm text-red-600">{errors.confirmPassword.message}</p>
              )}
            </div>
          </div>

          <div className="flex items-center justify-end gap-4 pt-2">
            <Link
              to="/login"
              className="text-base text-[#8d8d8d] transition-colors hover:text-[#606060]"
            >
              Back
            </Link>
            <button
              type="submit"
              disabled={isLoading}
              className="h-11 min-w-24 rounded-md bg-black px-6 text-base font-medium text-white transition-colors hover:bg-[#1a1a1a] disabled:cursor-not-allowed disabled:opacity-60"
            >
              {isLoading ? 'Creating...' : 'Create'}
            </button>
          </div>

          <p className="pt-1 text-center text-sm text-[#8d8d8d]">
            Already have an account?{' '}
            <Link to="/login" className="font-medium text-[#1f2022] hover:underline">
              Log in
            </Link>
          </p>
        </form>
      </div>
    </div>
  )
}
