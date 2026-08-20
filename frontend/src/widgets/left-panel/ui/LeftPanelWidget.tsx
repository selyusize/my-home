import { ProfileMiniCard } from '@entities/profile'

const LeftPanelWidget = () => {
  return (
    <aside className="w-[325px] border">
      <ProfileMiniCard />
    </aside>
  )
}

export { LeftPanelWidget }
