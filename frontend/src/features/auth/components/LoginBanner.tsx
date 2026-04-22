/**
 * Login banner component - Enhanced left side illustration
 */
export function LoginBanner() {
  return (
    <div className="banner-wrap">
      {/* Animated background decorations */}
      <div className="banner-bg-decoration-1" />
      <div className="banner-bg-decoration-2" />
      <div className="banner-bg-decoration-3" />

      {/* Brand illustration */}
      <img
        src="/favicon.svg"
        alt="Aitestos Platform"
        className="banner-image"
      />

      {/* Slogan */}
      <div className="banner-slogan">因为热爱 快乐成长</div>
    </div>
  )
}
