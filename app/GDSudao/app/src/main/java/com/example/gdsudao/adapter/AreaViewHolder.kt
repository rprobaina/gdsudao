package com.example.gdsudao.adapter

import android.view.LayoutInflater
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView
import com.example.gdsudao.R
import com.example.gdsudao.model.Area

class AreaViewHolder(inflater: LayoutInflater, parent: ViewGroup): RecyclerView.ViewHolder(inflater.inflate(R.layout.adapter_area, parent, false)){
    private var tvNome: TextView? = null
    private var tvDataPlantio: TextView? = null
    private var tvProxCorte: TextView? = null
    private var tvNumeroCortes: TextView? = null
    private var tvProgresso: TextView? = null
    //private var ivIcone: ImageView? = null
    private var tvDiarios: TextView? = null
    private var tvPrevisoes: TextView? = null
    private var tvNormais: TextView? = null


    init {
        tvNome = itemView.findViewById(R.id.tvNomeArea)
        tvDataPlantio = itemView.findViewById(R.id.tvDataPlantio)
        tvProxCorte = itemView.findViewById(R.id.tvDataProxCorte)
        tvNumeroCortes = itemView.findViewById(R.id.tvNumeroCortes)
        tvProgresso = itemView.findViewById(R.id.tvProgresso)
        tvDiarios = itemView.findViewById(R.id.tvDiarios)
        tvPrevisoes = itemView.findViewById(R.id.tvPrevisoes)
        tvNormais = itemView.findViewById(R.id.tvNormais)
        //ivIcone = itemView.findViewById(R.id.ivIconeSudao)
    }

    fun bind(area: Area){
        val ST_PRIRO_CORTE = 358.0f
        val ST_OUTROS_CORTES = 281.0f



        tvNome?.text = area.nome
        tvDataPlantio?.text = "Data Plantio: " + area.dataCorte
        tvProxCorte?.text = area.proxcorte.substring(0, 10).replace("-","/", false)
        tvNumeroCortes?.text = "Número de cortes: " + area.numeroCorte

        if (area.diario.length > 5){
            tvDiarios?.text = area.diario.substring(0, 4) + "%" //.substring(5) + "%"
        }else{
            tvDiarios?.text = area.diario + "%"
        }

        if (area.previsao.length > 5){
            tvPrevisoes?.text = area.previsao.substring(0, 4) + "%"  //.substring(5) + "%"
        }else{
            tvPrevisoes?.text = area.previsao + "%"  //.substring(5) + "%"
        }

        if (area.normal.length > 5){
            tvNormais?.text = area.normal.substring(0, 4) + "%"  //.substring(5) + "%"
        }else{
            tvNormais?.text = area.normal + "%"  //.substring(5) + "%"
        }

        // TODO: calcular o progresso com base nas constatnes de st e no numero de cortes
        // tvProgresso?.text = "Data Plantio: " + "|d: " + area.diario + "|p: " + area.previsao +  "|n: " + area.normal

        var progresso: Float
        if (area.numeroCorte.toInt() < 1) {
            progresso = (area.st.toFloat() / ST_PRIRO_CORTE) * 100
        }else {
            progresso =  (area.st.toFloat() / ST_OUTROS_CORTES ) * 100
        }

        tvProgresso?.text = "Progresso: " + progresso.toString().substring(0, 4) + "%"


        /*
        TODO: Implementar o esquema da figura ???
         */
    }



}

