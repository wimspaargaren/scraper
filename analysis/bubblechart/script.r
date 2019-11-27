library(ggplot2)
data <- read.csv(file="input.csv", header=TRUE, sep=",")
# maxVal <- sqrt(7 /pi)
# minVal <- sqrt(1 /pi)

# data$radius <- max(minVal,min(maxVal, sqrt( data$amount / pi  ) ) )
data$radius <- sqrt( data$amount / pi  )
 
print(data$radius)
# Could also be stored as pdf or svg
ggsave("result.png",
ggplot(data,aes(fintech_topic,ai_type))+
  geom_point(aes(size=radius*10),shape=21,fill="white")+
  geom_text(aes(label=amount),size=4)+
  # geom_text(aes(label=stat,vjust=radius*0.5+0.5,hjust=radius*-0.4),size=4)+
  scale_size_identity()+
  # Theme docu: https://www.datanovia.com/en/blog/ggplot-theme-background-color-and-grids/s
  theme(panel.grid.major=element_line(linetype=2,color="black"),
        panel.background = element_rect(fill="white"),
        axis.text.x=element_text(angle=90,hjust=1,vjust=0,color="black"),
        axis.text.y=element_text(angle=0,hjust=0.5,vjust=0,color="black")
        )
        )

