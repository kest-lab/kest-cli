import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { toast } from 'sonner'
import { useAuthStore, setAuthTokens } from '@/store/auth-store'
import { authApi } from '@/services/auth'

const loginSchema = z.object({
  username: z.string().min(1, 'Username is required'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
})

type LoginFormData = z.infer<typeof loginSchema>

export function LoginPage() {
  const navigate = useNavigate()
  const { setUser } = useAuthStore()
  const [isLoading, setIsLoading] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  })

  const onSubmit = async (data: LoginFormData) => {
    setIsLoading(true)
    try {
      const { user, token } = await authApi.login({
        username: data.username,
        password: data.password,
      })
      
      setAuthTokens(token)
      setUser(user)
      toast.success('Login successful!')
      navigate('/projects')
    } catch (error: any) {
      console.error('Login error:', error)
      const msg = error?.message || 'Login failed';
      toast.error(msg);
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-[#efefef] px-4 py-8 sm:px-6">
      <div className="mx-auto flex min-h-[calc(100vh-4rem)] w-full max-w-md items-center justify-center">
        <form onSubmit={handleSubmit(onSubmit)} className="w-full space-y-5">
          <div>
            <h1 className="text-[2.1rem] font-semibold tracking-tight text-[#1c1d20]">
              Log in to Kest
            </h1>
          </div>

          <div className="space-y-4">
            <div className="space-y-2.5">
              <Label htmlFor="username" className="text-[1.05rem] font-medium text-[#1f2022]">
                Email
              </Label>
              <Input
                id="username"
                type="text"
                autoComplete="username"
                placeholder="me@somewhere.com"
                {...register('username')}
                className="h-12 rounded-md border-[#d7d7d7] bg-[#ececec] px-4 text-base text-[#1f2022] placeholder:text-[#a0a0a0] focus-visible:ring-0"
              />
              {errors.username && (
                <p className="text-sm text-red-600">{errors.username.message}</p>
              )}
            </div>

            <div className="space-y-2.5">
              <Label htmlFor="password" className="text-[1.05rem] font-medium text-[#1f2022]">
                Password
              </Label>
              <Input
                id="password"
                type="password"
                autoComplete="current-password"
                placeholder="••••••••••"
                {...register('password')}
                className="h-12 rounded-md border-[#d7d7d7] bg-[#ececec] px-4 text-base text-[#1f2022] placeholder:text-[#8f8f8f] focus-visible:ring-0"
              />
              {errors.password && (
                <p className="text-sm text-red-600">{errors.password.message}</p>
              )}
            </div>
          </div>

          <div className="flex items-center justify-end gap-4 pt-2">
            <Link
              to="/"
              className="text-base text-[#8d8d8d] transition-colors hover:text-[#606060]"
            >
              Back
            </Link>
            <button
              type="submit"
              disabled={isLoading}
              className="h-11 min-w-24 rounded-md bg-black px-6 text-base font-medium text-white transition-colors hover:bg-[#1a1a1a] disabled:cursor-not-allowed disabled:opacity-60"
            >
              {isLoading ? 'Log in...' : 'Log in'}
            </button>
          </div>

          <p className="pt-1 text-center text-sm text-[#8d8d8d]">
            No account?{' '}
            <Link to="/register" className="font-medium text-[#1f2022] hover:underline">
              Register
            </Link>
          </p>
        </form>
      </div>
    </div>
  )
}
