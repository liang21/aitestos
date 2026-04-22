/**
 * Login banner component - Enhanced left side illustration
 */
export function LoginBanner() {
  return (
    <div className="w-1/2 h-screen flex items-start justify-start bg-gradient-to-br from-indigo-500 via-purple-500 to-purple-600 relative overflow-hidden m-0 p-0">
      {/* Animated background decorations */}
      <div className="absolute w-[500px] h-[500px] rounded-full bg-white/15 top-[-150px] left-[-150px] animate-float" />
      <div className="absolute w-[400px] h-[400px] rounded-full bg-white/10 bottom-[-100px] right-[-100px] animate-float-reverse" />
      <div className="absolute w-[200px] h-[200px] rounded-full bg-white/8 top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 animate-pulse-slow" />

      {/* Brand illustration */}
      <img
        src="/favicon.svg"
        alt="Aitestos Platform"
        className="absolute inset-0 w-full h-full object-cover opacity-95 m-0 p-0"
      />

      {/* Slogan */}
      <div className="absolute top-15 left-15 text-white/95 text-3xl font-light tracking-widest z-10 animate-fade-in-up">
        因为热爱 快乐成长
      </div>
    </div>
  )
}
